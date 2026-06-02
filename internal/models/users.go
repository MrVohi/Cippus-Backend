package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string   `gorm:"uniqueIndex;not null"`
	Email        string   `gorm:"uniqueIndex;not null" json:"-"`
	Bio          string
	PasswordHash string `json:"-"`
	AvatarURL    string
	Role         UserRole
}
