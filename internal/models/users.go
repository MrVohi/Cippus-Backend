package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string 
	Email        string
	Bio          string
	PasswordHash string `json:"-"`
	AvatarURL    string
	Role         UserRole
}
