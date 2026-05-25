package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	UserId      uint
	Author      User `gorm:"foreignKey:UserID"`
	Title       string
	Description string
}
