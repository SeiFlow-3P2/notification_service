package domain

import "time"

// NotificationRequest представляет структуру сообщения из Kafka
type NotificationRequest struct {
	EventID   string    `json:"event_id"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	UserID    int64     `json:"user_id"`
}