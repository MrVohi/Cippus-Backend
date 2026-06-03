package services

import (
	"cippus-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) SendMessage(senderID uint, receiverID uint, content string) (*models.Message, error) {
	msg := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
	result := s.db.Create(&msg)
	return &msg, result.Error
}

func (s *MessageService) GetConversations(userID uint) ([]models.Message, error) {
	var messages []models.Message
	s.db.Raw(`SELECT DISTINCT ON (LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id)) * FROM messages WHERE sender_id = ? OR receiver_id = ? ORDER BY LEAST(sender_id, receiver_id), GREATEST(sender_id, receiver_id), created_at DESC`, userID, userID).Scan(&messages)
	return messages, nil
}

func (s *MessageService) GetThread(userID uint, otherUserID uint) ([]models.Message, error) {
	var messages []models.Message
	s.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, otherUserID, otherUserID, userID).Order("created_at ASC").Find(&messages)
	return messages, nil
}

func (s *MessageService) MarkAsRead(messageID uint, userID uint) error {
	result := s.db.Model(&models.Message{}).Where("id = ? AND receiver_id = ?", messageID, userID).Update("read_at", time.Now())
	return result.Error
}
