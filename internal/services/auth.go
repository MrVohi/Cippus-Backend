package services

import (
	"cippus-backend/internal/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

func (s *AuthService) Refresh(token string) (string, string, error) {
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

func (s *AuthService) Logout(userID uint) error {
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

func (s *AuthService) PasswordResetRequest(email string) {
	user := models.User{}
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return
	}

	raw := make([]byte, 32)
	rand.Read(raw)
	token := hex.EncodeToString(raw)

	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pst := models.PasswordResetToken{
		UserID:    user.ID,
		User:      user,
		TokenHash: string(hash),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	result = s.db.Create(&pst)

	client := resend.NewClient(s.resendApiKey)

	params := resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{user.Email},
		Subject: "Password reset",
		Html:    "<p>Reset link: http://localhost:3000/auth/password-reset?token=" + token + "</p>",
	}

	client.Emails.Send(&params)
}

func (s *AuthService) PasswordResetConfirm(token string, newPassword string) error {
	pst := models.PasswordResetToken{}
	result := s.db.Where("used_at is NULL").Where("expires_at > ?", time.Now()).First(&pst)
	if result.Error != nil {
		return fmt.Errorf("Invalid or expired token")
	}

	err := bcrypt.CompareHashAndPassword([]byte(pst.TokenHash), []byte(token))
	if err != nil {
		return fmt.Errorf("Invalid or expired password")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Cannot generate hash for password")
	}

	s.db.Model(&models.User{}).Where("id = ?", pst.UserID).Update("password_hash", string(hash))

	now := time.Now()
	s.db.Model(&pst).Update("used_at", &now)
	return nil
}
