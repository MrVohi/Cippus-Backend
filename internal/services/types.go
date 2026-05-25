package services

import (
	"cippus-backend/internal/models"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type AuthService struct {
	db              *gorm.DB
	secret          string
	resendApiKey    string
	recaptchaSecret string
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	User         models.User
}

type UserService struct {
	db *gorm.DB
}

type PostService struct {
	db *gorm.DB
}

type PostInput struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	ImageURL    string `json:"imageUrl"`
	CategoryIDs []uint `json:"categoryIds"`
	Stuck       *bool  `json:"stuck"`
}

type GetPostsFilter struct {
	CategoryID *uint
	AuthorID   *uint
	Stuck      *bool
	Q          *string
}

type MinioService struct {
	Client   *minio.Client
	Bucket   string
	Endpoint string
}

type CategoryService struct {
	db *gorm.DB
}

type CategoryInput struct {
	Name        string
	Description string
}
