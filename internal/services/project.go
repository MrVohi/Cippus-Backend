package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func (s *ProjectService) CreateProject(UserId uint, input ProjectInput) (*models.Project, error) {
	project := &models.Project{
		UserId:      UserId,
		Title:       input.Title,
		Description: input.Description,
	}

	result := s.db.Create(&project)
	if result.Error != nil {
		return &models.Project{}, fmt.Errorf("Could not save project to db")
	}

	return project, nil
}

func (s *ProjectService) GetProjectsByUser(UserId uint) ([]models.Project, error) {
	projects := []models.Project{}

	result := s.db.Where("user_id = ?", UserId).Find(&projects)
	if result.Error != nil {
		return []models.Project{}, fmt.Errorf("could not find projects")
	}

	return projects, nil
}

func (s *ProjectService) GetProjectById(id uint) (models.Project, error) {
	projects := models.Project{}

	result := s.db.Where("id = ?", id).First(&projects)
	if result.Error == gorm.ErrRecordNotFound {
		return models.Project{}, fmt.Errorf("could not find project")
	} else if result.Error != nil {
		return models.Project{}, fmt.Errorf("error with db")
	}

	return projects, nil
}

func (s *ProjectService) UpdateProject(id uint, callerID uint, callerRole string, input ProjectInput) (models.Project, error) {
	project := models.Project{}

	result := s.db.Where("id = ?", id).First(&project)
	if result.Error != nil {
		return models.Project{}, fmt.Errorf("could not find project")
	}

	if callerID != project.UserId && callerRole == "user" {
		return models.Project{}, fmt.Errorf("forbidden")
	}

	project.Title = input.Title
	project.Description = input.Description

	result = s.db.Save(&project)
	if result.Error != nil {
		return models.Project{}, fmt.Errorf("Could not update project")
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(id uint, callerID uint, callerRole string) error {
	project := models.Project{}

	result := s.db.Where("id = ?", id).First(&project)
	if result.Error != nil {
		return fmt.Errorf("could not find project")
	}

	if callerID != project.UserId && callerRole == "user" {
		return fmt.Errorf("forbidden")
	}

	result = s.db.Delete(&project)
	if result.Error != nil {
		return fmt.Errorf("Could not delete project")
	}

	return nil
}

func (s *ProjectService) AddPostToProject(projectID uint, postID uint, callerID uint) error {
	project := models.Project{}

	result := s.db.Where("id = ?", projectID).First(&project)
	if result.Error != nil {
		return fmt.Errorf("could not find project")
	}

	if callerID != project.UserId {
		return fmt.Errorf("forbidden")
	}

	post := models.Post{}

	result = s.db.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return fmt.Errorf("could not find post")
	}

	if callerID != post.UserId {
		return fmt.Errorf("forbidden")
	}

	projectPost := models.ProjectPost{
		ProjectID: projectID,
		PostID:    postID,
	}

	result = s.db.Create(&projectPost)
	if result.Error != nil {
		return fmt.Errorf("Could not add post to project")
	}

	return nil
}

func (s *ProjectService) RemovePostFromProject(projectID uint, postID uint, callerID uint) error {
	project := models.Project{}

	result := s.db.Where("id = ?", projectID).First(&project)
	if result.Error != nil {
		return fmt.Errorf("could not find project")
	}

	if callerID != project.UserId {
		return fmt.Errorf("forbidden")
	}

	post := models.Post{}

	result = s.db.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return fmt.Errorf("could not find post")
	}

	if callerID != post.UserId {
		return fmt.Errorf("forbidden")
	}

	projectPost := models.ProjectPost{
		ProjectID: projectID,
		PostID:    postID,
	}

	result = s.db.Delete(&projectPost)
	if result.Error != nil {
		return fmt.Errorf("Could not remove post from project")
	}

	return nil
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{db: db}
}
