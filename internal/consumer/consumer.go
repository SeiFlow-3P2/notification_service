package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	"notification_service/domain"
	"NOTIFICATION_SERVICE/telegram"
)

type Consumer struct {
	reader    *kafka.Reader
	tgClient  *telegram.Client
	stopChan  chan struct{}
}

func NewConsumer(kafkaClient *kafka.Reader, tgClient *telegram.Client, topic string) (*Consumer, error) {
	kafkaClient.Config().Topic = topic
	return &Consumer{
		reader:    kafkaClient,
		tgClient:  tgClient,
		stopChan:  make(chan struct{}),
	}, nil
}

func (c *Consumer) Start() {
	for {
		select {
		case <-c.stopChan:
			log.Println("Consumer shutting down...")
			return
		default:
			msg, err := c.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Failed to read message from Kafka: %v", err)
				continue
			}

			var req domain.NotificationRequest
			if err := json.Unmarshal(msg.Value, &req); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if req.EventID == "" || req.Title == "" || req.UserID == 0 {
				log.Printf("Invalid message: missing required fields")
				continue
			}

			messageText := fmt.Sprintf("Событие началось: %s", req.Title)
			if err := c.sendWithRetry(req.UserID, messageText); err != nil {
				log.Printf("Failed to send message to user %d: %v", req.UserID, err)
				continue
			}
			log.Printf("Notification sent to user %d: %s", req.UserID, messageText)
		}
	}
}

func (c *Consumer) sendWithRetry(userID int64, text string) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 5 * time.Minute
	return backoff.Retry(func() error {
		return c.tgClient.SendMessage(userID, text)
	}, b)
}

func (c *Consumer) Stop() {
	close(c.stopChan)
	c.reader.Close()
}