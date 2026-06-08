package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

func (s *SearchService) SemanticSearch(query string, limit int) ([]models.Post, error) {
	embed, err := s.es.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	result := s.db.Raw("SELECT *, embedding <=> ? AS distance FROM posts WHERE deleted_at IS NULL ORDER BY distance ASC LIMIT ?", pgvector.NewVector(embed), limit).Scan(&posts)
	if result.Error != nil {
		return nil, result.Error
	}

	return posts, nil
}

func (s *SearchService) GetSimilar(postID uint, limit int) ([]models.Post, error) {
	var post models.Post
	result := s.db.Select("embedding").First(&post, postID)
	if result.Error != nil {
		return []models.Post{}, result.Error
	}

	if len(post.Embedding.Slice()) == 0 {
		return []models.Post{}, fmt.Errorf("no embedding")
	}
	var similar []models.Post
	result = s.db.Raw("SELECT * FROM posts WHERE deleted_at IS NULL AND id != ? ORDER BY embedding <=> ? ASC LIMIT ?", postID, post.Embedding, limit).Scan(&similar)
	if result.Error != nil {
		return []models.Post{}, result.Error
	}

	return similar, nil
}

func NewSearchService(db *gorm.DB, es *EmbeddingService) *SearchService {
	return &SearchService{db: db, es: *es}
}
