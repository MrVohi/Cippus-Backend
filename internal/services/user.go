package services

import (
	"cippus-backend/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func (s *UserService) PatchMe(userID uint, fields map[string]interface{}) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(fields)
	return result.Error
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
