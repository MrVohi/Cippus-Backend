package services

import (
	"cippus-backend/internal/models"

	"github.com/minio/minio-go/v7"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type AuthService struct {
	db              *gorm.DB
	secret          string
	frontendURL     string
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
	es EmbeddingService
}

type PostInput struct {
	Title       string `json:"title"       binding:"required,min=3,max=200"`
	Content     string `json:"content"     binding:"required,min=10"`
	CategoryIDs []uint `json:"categoryIds" binding:"required,min=1"`
	Stuck       *bool  `json:"stuck"       binding:"omitempty"`
}

type GetPostsFilter struct {
	CategoryID *uint
	AuthorID   *uint
	Stuck      *bool
	Q          *string
	Sort       string
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
	Name        string `json:"name"        binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"omitempty,max=255"`
}

type ProjectService struct {
	db *gorm.DB
}

type ProjectInput struct {
	Title       string `json:"title"       binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type MessageService struct {
	db *gorm.DB
}

type PushService struct {
	db              *gorm.DB
	vapidPublic     string
	vapidPrivate    string
	vapidSubscriber string
}

type RabbitPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

type NotificationService struct {
	db     *gorm.DB
	rabbit *RabbitPublisher
}

type EmbeddingService struct {
	url        string
	embedModel string
	model      string
}

type SearchService struct {
	db *gorm.DB
	es EmbeddingService
}

type CommentService struct {
	db *gorm.DB
}