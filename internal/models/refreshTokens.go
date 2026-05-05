package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint
	TokenHash string
	ExpiresAt time.Time
}
