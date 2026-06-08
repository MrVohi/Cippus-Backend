package services

import (
	"cippus-backend/internal/models"
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"gorm.io/gorm"
)

func NewPushService(db *gorm.DB, public string, private string, subscriber string) *PushService {
	return &PushService{db: db, vapidPublic: public, vapidPrivate: private, vapidSubscriber: subscriber}
}

func (s *PushService) Subscribe(userID uint, endpoint string, p256dh string, auth string) error {
	sub := models.PushSubscription{
		UserID:    userID,
		Endpoint:  endpoint,
		P256dhKey: p256dh,
		AuthKey:   auth,
	}
	return s.db.Create(&sub).Error
}

func (s *PushService) Unsubscribe(userID uint, endpoint string) error {
	return s.db.Where("user_id = ? AND endpoint = ?", userID, endpoint).Delete(&models.PushSubscription{}).Error
}

func (s *PushService) SendPush(userID uint, title string, body string) error {

	var subs []models.PushSubscription
	s.db.Where("user_id = ?", userID).Find(&subs)

	payload, _ := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})

	for _, sub := range subs {
		subscription := &webpush.Subscription{Endpoint: sub.Endpoint, Keys: webpush.Keys{
				P256dh: sub.P256dhKey,
				Auth:   sub.AuthKey,
			},
		}
		webpush.SendNotification(payload, subscription, &webpush.Options{
			VAPIDPublicKey:  s.vapidPublic,
			VAPIDPrivateKey: s.vapidPrivate,
			Subscriber:      s.vapidSubscriber,
			TTL:             30,
		})
	}

	return nil
}
