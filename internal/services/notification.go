package services

import (
	"cippus-backend/internal/models"
	"encoding/json"

	"gorm.io/gorm"
)

func NewNotificationService(db *gorm.DB, rabbit *RabbitPublisher) *NotificationService {
	return &NotificationService{db: db, rabbit: rabbit}
}

func (s *NotificationService) Notify(
	recipientID uint,
	actorID uint,
	notifType models.NotificationType,
	entityType models.ContentType,
	entityID uint,
) error {
	notif := models.Notification{
		RecipientID: recipientID,
		ActorID:     actorID,
		Type:        notifType,
		EntityType:  entityType,
		EntityID:    entityID,
	}
	if err := s.db.Create(&notif).Error; err != nil {
		return err
	}

	payload, _ := json.Marshal(map[string]any{
		"recipientID": recipientID,
		"actorID":     actorID,
		"type":        notifType,
		"entityType":  entityType,
		"entityID":    entityID,
	})
	s.rabbit.Publish("notifications.dispatch", payload)

	return nil
}
