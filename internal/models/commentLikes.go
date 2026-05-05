package models

import "gorm.io/gorm"

type CommentLike struct {
	gorm.Model
	UserID    uint
	User      User
	CommentID uint
	Comment   Comment
	Liked     bool
}
