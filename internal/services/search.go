package services

import (
	"cippus-backend/internal/models"
	"fmt"
	"sort"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type PostWithDistance struct {
	Post     models.Post
	Distance float64
}

func (s *SearchService) SemanticSearch(query string, limit int) ([]PostWithDistance, error) {
	embed, err := s.es.GenerateEmbedding(query)
	if err != nil {
		return nil, err
	}

	var rows []struct {
		ID       uint
		Distance float64
	}
	result := s.db.Raw("SELECT id, embedding <=> ? AS distance FROM posts WHERE deleted_at IS NULL ORDER BY distance ASC LIMIT ?", pgvector.NewVector(embed), limit).Scan(&rows)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(rows) == 0 {
		return []PostWithDistance{}, nil
	}

	ids := make([]uint, len(rows))
	distByID := make(map[uint]float64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
		distByID[r.ID] = r.Distance
	}

	var posts []models.Post
	if err := s.db.Preload("Author").Preload("Categories").Where("id IN ?", ids).Find(&posts).Error; err != nil {
		return nil, err
	}

	out := make([]PostWithDistance, len(posts))
	for i, p := range posts {
		out[i] = PostWithDistance{Post: p, Distance: distByID[p.ID]}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Distance < out[j].Distance
	})

	return out, nil
}

func (s *SearchService) GetSimilar(postID uint, limit int) ([]PostWithDistance, error) {
	var post models.Post
	result := s.db.Select("embedding").First(&post, postID)
	if result.Error != nil {
		return []PostWithDistance{}, result.Error
	}

	if len(post.Embedding.Slice()) == 0 {
		return []PostWithDistance{}, fmt.Errorf("no embedding")
	}

	var rows []struct {
		ID       uint
		Distance float64
	}
	result = s.db.Raw("SELECT id, embedding <=> ? AS distance FROM posts WHERE deleted_at IS NULL AND id != ? AND embedding <=> ? < 0.3 ORDER BY distance ASC LIMIT ?", post.Embedding, postID, post.Embedding, limit).Scan(&rows)
	if result.Error != nil {
		return []PostWithDistance{}, result.Error
	}
	if len(rows) == 0 {
		return []PostWithDistance{}, nil
	}

	ids := make([]uint, len(rows))
	distByID := make(map[uint]float64, len(rows))
	for i, r := range rows {
		ids[i] = r.ID
		distByID[r.ID] = r.Distance
	}

	var similar []models.Post
	if err := s.db.Preload("Author").Preload("Categories").Where("id IN ?", ids).Find(&similar).Error; err != nil {
		return []PostWithDistance{}, err
	}

	out := make([]PostWithDistance, len(similar))
	for i, p := range similar {
		out[i] = PostWithDistance{Post: p, Distance: distByID[p.ID]}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Distance < out[j].Distance
	})

	return out, nil
}

func NewSearchService(db *gorm.DB, es *EmbeddingService) *SearchService {
	return &SearchService{db: db, es: *es}
}
