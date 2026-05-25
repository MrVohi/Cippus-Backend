package models

type ProjectPost struct {
	ProjectID uint `gorm:"primaryKey"`
	Project   Project
	PostID    uint `gorm:"primaryKey"`
	Post      Post
}
