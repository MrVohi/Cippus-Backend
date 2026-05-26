package services

import (
	"cippus-backend/internal/models"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (s *AuthService) Register(email string, username string, password string) (AuthResult, error) {
	row := models.User{}

	usernameRow := models.User{}
	usernameResult := s.db.Where("username = ?", username).First(&usernameRow)
	if usernameResult.Error == nil {
		suggestedUsername := GenerateUniqueUsername(s.db, username)
		return AuthResult{}, fmt.Errorf("username already exists: %s", suggestedUsername)
	}

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

func GenerateUniqueUsername(db *gorm.DB, base string) string {

	for {
		raw := make([]byte, 2)
		rand.Read(raw)

		var candidate string

		suffix := (int(raw[0])<<8|int(raw[1]))%9000 + 1000
		candidate = fmt.Sprintf("%s%d", base, suffix)

		row := models.User{}
		result := db.Where("username = ?", candidate).First(&row)
		if result.Error == gorm.ErrRecordNotFound {
			return candidate
		}
	}
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

func (s *AuthService) Refresh(token string) (string, string, models.User, error) {
	userID, err := ParseRefreshToken(token, s.secret)
	if err != nil {
		return "", "", models.User{}, fmt.Errorf("Cannot find userID")
	}

	user := models.User{}
	result := s.db.First(&user, userID)
	if result.Error != nil {
		return "", "", models.User{}, fmt.Errorf("Cannot find User")
	}

	newRefreshToken, _, err := RotateRefreshToken(s.db, userID, []byte(token), s.secret)
	if err != nil {
		return "", "", models.User{}, fmt.Errorf("Cannot rotate refresh token")
	}

	newAccessToken, err := GenerateAccessToken(userID, user.Role, s.secret)
	if err != nil {
		return "", "", models.User{}, fmt.Errorf("Cannot generate new access token")
	}

	return newRefreshToken, newAccessToken, user, nil
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

func (s *AuthService) PasswordResetRequest(email string) (error){
	user := models.User{}
	result := s.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return fmt.Errorf("Cannot find user")
	}

	raw := make([]byte, 32)
	rand.Read(raw)
	token := hex.EncodeToString(raw)

	hash := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(hash[:])

	pst := models.PasswordResetToken{
		UserID:    user.ID,
		User:      user,
		TokenHash: string(hashHex),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	result = s.db.Create(&pst)
	if result.Error != nil {
		return fmt.Errorf("Cannot create password reset token")
	}

	client := resend.NewClient(s.resendApiKey)

	params := resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{user.Email},
		Subject: "Password reset",
		Html:    "<p>Reset link: http://localhost:3000/auth/password-reset?token=" + token + "</p>",
	}

	client.Emails.Send(&params)
	return nil
}

func (s *AuthService) PasswordResetConfirm(token string, newPassword string) error {
	var tokens []models.PasswordResetToken
	result := s.db.Where("used_at is NULL").Where("expires_at > ?", time.Now()).Find(&tokens)
	if result.Error != nil {
		return fmt.Errorf("Invalid or expired token")
	}

	incomingHash := sha256.Sum256([]byte(token))
	incomingHashHex := hex.EncodeToString(incomingHash[:])

	var found *models.PasswordResetToken

	for _, pst := range tokens {
		if pst.TokenHash == incomingHashHex {
			found = &pst
			break
		}
	}

	if found == nil {
		return fmt.Errorf("invalid or expired token")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Cannot generate hash for password")
	}

	result = s.db.Model(&models.User{}).Where("id = ?", found.UserID).Update("password_hash", string(hash))
	if result.Error != nil {
		return fmt.Errorf("Cannot update password")
	}

	now := time.Now()
	result = s.db.Model(found).Update("used_at", &now)
	if result.Error != nil {
		return fmt.Errorf("Cannot update password")
	}
	return nil
}

func NewAuthService(db *gorm.DB, secret, resendApiKey, recaptchaSecret string) *AuthService {
	return &AuthService{db: db, secret: secret, resendApiKey: resendApiKey, recaptchaSecret: recaptchaSecret}
}
