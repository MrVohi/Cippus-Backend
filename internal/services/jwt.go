package services

import (
	"cippus-backend/internal/models"
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID uint, role models.UserRole, secret string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"userRole": role,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))

	return token, err
}

func ValidateAccessToken(token string, secret string) (jwt.MapClaims, error){
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if (!ok) {
			return nil, errors.New("Unexpected signing method:")
		}
		return []byte(secret), nil
	})
	return claims, err
}