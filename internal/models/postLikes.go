package models

import "gorm.io/gorm"

type PostLike struct {
	gorm.Model
	UserID uint
	User   User
	PostID uint
	Post   Post
	Liked  bool
}
