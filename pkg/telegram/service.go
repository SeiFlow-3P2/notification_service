package telegram

import (
    "context"
    "fmt"

    "notification_service/client/telegram"
)

type Service struct {
    client *telegram.Client
}

func NewService(client *telegram.Client) *Service {
    return &Service{client: client}
}

func (c *Service) SendNotification(ctx context.Context, userID int64, title string) error {
    return c.client.SendMessage(ctx, userID, fmt.Sprintf("Событие началось: %s", title))
}