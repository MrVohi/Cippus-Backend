package services

import (
	"cippus-backend/internal/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GenerateRefreshToken() ([]byte, []byte, error) {
	raw := make([]byte, 32)
	_, err := rand.Read(raw)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to generate 32 random bytes: %w", err)
	}

	token := make([]byte, hex.EncodedLen(len(raw)))
	hex.Encode(token, raw)

	hash, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("Unable to generate hash: %w", err)
	}

	return token, hash, nil
}

func StoreRefreshToken(db *gorm.DB, userID uint, hash []byte) error {
	pk := db.Create(&models.RefreshToken{
		UserID:    userID,
		TokenHash: string(hash),
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

func RotateRefreshToken(db *gorm.DB, userID uint, incoming []byte) ([]byte, []byte, error) {
	rt, err := FindRefreshToken(db, userID)
	if err != nil {
		return nil, nil, err
	}

	result := bcrypt.CompareHashAndPassword([]byte(rt.TokenHash), incoming)

	if result != nil {
		return nil, nil, fmt.Errorf("Refresh token for User %d does not match database! %w", userID, result)
	}
	err = DeleteRefreshToken(db, rt.ID)
	if err != nil {
		return nil, nil, err
	}

	newPlain, newHash, err := GenerateRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	err = StoreRefreshToken(db, userID, newHash)
	if err != nil {
		return nil, nil, err
	}

	return newPlain, newHash, nil

}
