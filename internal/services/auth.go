package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	secret string
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         models.User
}

func (s *AuthService) Register(email string, username string, password string) (AuthResult, error) {
	row := models.User{}

	result := s.db.Where("email = ?", email).First(&row)
	if result.Error == nil {
		return AuthResult{}, fmt.Errorf("email already exists")
	} else if result.Error != gorm.ErrRecordNotFound {
		return AuthResult{}, fmt.Errorf("error with db")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate hash")
	}
	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: string(hash),
		Role:         models.RoleUser,
	}
	result = s.db.Create(&user)
	if result.Error != nil {
		return AuthResult{}, fmt.Errorf("cannot register new user")
	}

	accessToken, err := GenerateAccessToken(user.ID, user.Role, s.secret)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate access token")
	}
	refreshToken, tokenHash, err := GenerateRefreshToken(user.ID, s.secret)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate refresh token")
	}
	err = StoreRefreshToken(s.db, user.ID, tokenHash)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot save refresh token")
	}

	return AuthResult{accessToken, refreshToken, user}, nil
}

func (s *AuthService) Login(email string, password string) (AuthResult, error) {
	row := models.User{}
	result := s.db.Where("email = ?", email).First(&row)
	if result.Error == gorm.ErrRecordNotFound {
		return AuthResult{}, fmt.Errorf("invalid credentials")
	} else if result.Error != nil {
		return AuthResult{}, fmt.Errorf("error with db")
	} 

	if bcrypt.CompareHashAndPassword([]byte(row.PasswordHash), []byte(password)) != nil {
		return AuthResult{}, fmt.Errorf("invalid credentials")
	}

	accessToken, err := GenerateAccessToken(row.ID, row.Role, s.secret)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate access token")
	}
	refreshToken, tokenHash, err := GenerateRefreshToken(row.ID, s.secret)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate refresh token")
	}
	err = StoreRefreshToken(s.db, row.ID, tokenHash)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot save refresh token")
	}

	return AuthResult{accessToken, refreshToken, row}, nil
}

func (s *AuthService) Refresh(token string) (string, string, error){
	userID, err := ParseRefreshToken(token, s.secret)
	if err != nil {
		return "", "", fmt.Errorf("Cannot find userID")
	}

	user := models.User{}
	result := s.db.First(&user, userID)
	if result.Error != nil {
		return "", "", fmt.Errorf("Cannot find User")
	}

	newRefreshToken, _, err := RotateRefreshToken(s.db, userID, []byte(token), s.secret)
	if err != nil {
		return "", "", fmt.Errorf("Cannot rotate refresh token")
	}

	newAccessToken, err := GenerateAccessToken(userID, user.Role, s.secret)
	if err != nil {
		return "", "", fmt.Errorf("Cannot generate new access token")
	}

	return newRefreshToken, newAccessToken, nil
}

func (s *AuthService) Logout(userID uint) error{
	result, err := FindRefreshToken(s.db, userID)
	if err != nil {
		return fmt.Errorf("Cannot find refresh token for Logout")
	}
	err = DeleteRefreshToken(s.db, result.ID)
	if err != nil {
		return fmt.Errorf("Cannot delete refresh token for Logout")
	}

	return nil
}