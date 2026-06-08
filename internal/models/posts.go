package models

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserId           uint
	Author           User `gorm:"foreignKey:UserID"`
	Title            string
	Content          string `gorm:"type:text"`
	ImageURL         string
	Stuck            bool
	ModerationStatus ModerationStatus
	Categories       []Category      `gorm:"many2many:post_categories;"`
	Embedding        pgvector.Vector `gorm:"type:vector(768)"`
}
