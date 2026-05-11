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
	RefreshToken []byte
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
	refreshToken, tokenHash, err := GenerateRefreshToken()
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
	refreshToken, tokenHash, err := GenerateRefreshToken()
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot generate refresh token")
	}
	err = StoreRefreshToken(s.db, row.ID, tokenHash)
	if err != nil {
		return AuthResult{}, fmt.Errorf("cannot save refresh token")
	}

	return AuthResult{accessToken, refreshToken, row}, nil
}
