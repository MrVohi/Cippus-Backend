package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	rows := []models.Category{}
	result := s.db.Find(&rows)
	if result.Error != nil {
		return nil, result.Error
	}

	return rows, nil
}

func (s *CategoryService) CreateCategory(input CategoryInput) (models.Category, error) {
	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
	}

	result := s.db.Create(&category)
	if result.Error != nil {
		return models.Category{}, fmt.Errorf("cannot register new category")
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(cID uint, input CategoryInput) (models.Category, error) {
	category := models.Category{}
	result := s.db.First(&category, cID)
	if result.Error != nil {
		return models.Category{}, fmt.Errorf("category not found")
	}

	if input.Name != category.Name {
		category.Name = input.Name
	}
	if input.Description != category.Description {
		category.Description = input.Description
	}

	result = s.db.Save(&category)
	if result.Error != nil {
		return category, fmt.Errorf("Cannot save category to db")
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(cID uint) (error) {
	category := models.Category{}
	result := s.db.First(&category, cID)
	if result.Error != nil {
		return fmt.Errorf("category not found")
	}

	result = s.db.Delete(&category)
	if result.Error != nil {
		return fmt.Errorf("Cannot delete category")
	}
	
	return nil
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}