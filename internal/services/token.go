package services

import (
	"cippus-backend/internal/models"
	"encoding/hex"
	"fmt"
	"time"

	"crypto/sha256"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func GenerateRefreshToken(userID uint, secret string) (string, []byte, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().AddDate(0, 0, 7).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	hash := sha256.Sum256([]byte(token))

	return token, hash[:], nil
}

func StoreRefreshToken(db *gorm.DB, userID uint, hash []byte) error {
	pk := db.Create(&models.RefreshToken{
		UserID:    userID,
		TokenHash: hex.EncodeToString(hash),
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	})

	if pk.Error == nil {
		return nil
	} else {
		return fmt.Errorf("Failed to store Refresh token (%w) of User %d", pk.Error, userID)
	}
}

func FindRefreshToken(db *gorm.DB, userID uint) (models.RefreshToken, error) {
	row := models.RefreshToken{}
	result := db.Where("user_id = ?", userID).First(&row)

	return row, result.Error
}

func DeleteRefreshToken(db *gorm.DB, id uint) error {
	result := db.Delete(&models.RefreshToken{}, id)
	return result.Error
}

func RotateRefreshToken(db *gorm.DB, userID uint, incoming []byte, secret string) (string, []byte, error) {
	rt, err := FindRefreshToken(db, userID)
	if err != nil {
		return "", nil, err
	}

	hash := sha256.Sum256([]byte(incoming))

	if rt.TokenHash != hex.EncodeToString(hash[:]) {
		return "", nil, fmt.Errorf("Refresh token for User %d does not match database!", userID)
	}
	err = DeleteRefreshToken(db, rt.ID)
	if err != nil {
		return "", nil, err
	}

	newPlain, newHash, err := GenerateRefreshToken(userID, secret)
	if err != nil {
		return "", nil, err
	}

	err = StoreRefreshToken(db, userID, newHash)
	if err != nil {
		return "", nil, err
	}

	return newPlain, newHash, nil

}

func ParseRefreshToken(token string, secret string) (uint, error) {
	claims, err := ValidateAccessToken(token, secret)
	if err != nil {
		return 0, fmt.Errorf("Invalid access token")
	}

	raw := claims["userID"]
	userID, ok := raw.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid userID claim")
	}
	return uint(userID), nil
}
