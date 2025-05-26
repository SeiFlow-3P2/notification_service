package telegram

import (
	"fmt"
	"notification_service/client/telegram"
)

type Service struct {
	client *telegram.Client
}

func NewService(client *telegram.Client) *Service {
	return &Service{client: client}
}

func (s *Service) SendNotification(userID int64, title string) error {
	return s.client.SendMessage(userID, fmt.Sprintf("Событие началось: %s", title))
}