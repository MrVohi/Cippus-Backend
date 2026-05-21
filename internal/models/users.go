package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string
	Email        string
	Bio          string
	PasswordHash string
	AvatarURL    string
	Role         UserRole
}
