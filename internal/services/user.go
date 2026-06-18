package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func (s *UserService) PatchMe(userID uint, fields map[string]interface{}) error {
	result := s.db.Model(&models.User{}).Where("id = ?", userID).Updates(fields)
	return result.Error
}

func (s *UserService) GetUserByID(id uint) (models.User, error) {
	user := models.User{}
	result := s.db.Model(&models.User{}).Where("id = ?", id).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return models.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateAvatarURL(userID uint, url string) error {
	result := s.db.Model(models.User{}).Where("id = ?", userID).Update("avatar_url", url)
	return result.Error
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}
