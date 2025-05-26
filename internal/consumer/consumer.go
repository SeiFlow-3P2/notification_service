package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

    "github.com/SeiFlow-3P2/notification_service/internal/client/telegram"
	"github.com/SeiFlow-3P2/notification_service/internal/domain"
	"github.com/SeiFlow-3P2/notification_service/pkg/logger"
	// "github.com/SeiFlow-3P2/notification_service/pkg/telegram"
	"github.com/SeiFlow-3P2/notification_service/pkg/telemetry"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Consumer struct {
    reader    *kafka.Reader
    tgClient  *telegram.Client
    stopChan  chan struct{}
    tracer    trace.Tracer
}

func NewConsumer(kafkaClient *kafka.Reader, tgClient *telegram.Client, topic string) (*Consumer, error) {
    // kafkaTopic := kafkaClient.Config().Topic 
    // kafkaTopic = topic
    return &Consumer{
        reader:    kafkaClient,
        tgClient:  tgClient,
        stopChan:  make(chan struct{}),
        tracer:    telemetry.Tracer(),
    }, nil
}

func (c *Consumer) Start() {
    for {
        select {
        case <-c.stopChan:
            logger.Info("Consumer shutting down")
            return
        default:
            // Начало трассировки: чтение сообщения из Kafka
            ctx, span := c.tracer.Start(context.Background(), "kafka_read_message")
            defer span.End()

            msg, err := c.reader.ReadMessage(ctx)
            if err != nil {
                logger.Error("Failed to read message from Kafka",
                    err,
                    zap.String("topic", msg.Topic),
                    zap.Int("partition", msg.Partition),
                    zap.Int64("offset", msg.Offset))
                span.RecordError(err)
                continue
            }
            span.SetAttributes(
                attribute.String("kafka.topic", msg.Topic),
                attribute.Int("kafka.partition", msg.Partition),
                attribute.Int64("kafka.offset", msg.Offset),
            )

            // Трассировка: десериализация JSON
            ctx, unmarshalSpan := c.tracer.Start(ctx, "unmarshal_message")
            var req domain.NotificationRequest
            if err := json.Unmarshal(msg.Value, &req); err != nil {
                logger.Error("Failed to unmarshal message",
                    err,
                    zap.ByteString("message", msg.Value))
                unmarshalSpan.RecordError(err)
                unmarshalSpan.End()
                continue
            }
            unmarshalSpan.End()

            // Трассировка: валидация сообщения
            ctx, validateSpan := c.tracer.Start(ctx, "validate_message")
            if req.EventID == "" || req.Title == "" || req.UserID == 0 {
                logger.Warn("Invalid message: missing required fields",
                    zap.String("event_id", req.EventID),
                    zap.String("title", req.Title),
                    zap.Int64("user_id", req.UserID))
                validateSpan.SetAttributes(
                    attribute.String("validation.error", "missing required fields"),
                )
                validateSpan.End()
                continue
            }
            validateSpan.End()

            // Трассировка: отправка уведомления в Telegram
            messageText := fmt.Sprintf("Событие началось: %s", req.Title)
            ctx, sendSpan := c.tracer.Start(ctx, "send_telegram_message")
            if err := c.sendWithRetry(ctx, req.UserID, messageText); err != nil {
                logger.Error("Failed to send message to user",
                    err,
                    zap.Int64("user_id", req.UserID),
                    zap.String("message", messageText))
                sendSpan.RecordError(err)
                sendSpan.End()
                continue
            }
            sendSpan.SetAttributes(
                attribute.Int64("telegram.user_id", req.UserID),
                attribute.String("telegram.message", messageText),
            )
            sendSpan.End()

            logger.Info("Notification sent to user",
                zap.Int64("user_id", req.UserID),
                zap.String("message", messageText))
        }
    }
}

func (c *Consumer) sendWithRetry(ctx context.Context, userID int64, text string) error {
    b := backoff.NewExponentialBackOff()
    b.MaxElapsedTime = 5 * time.Minute
    return backoff.Retry(func() error {
        // Трассировка для каждой попытки отправки
        _, span := c.tracer.Start(ctx, "telegram_send_attempt")
        defer span.End()

        err := c.tgClient.SendMessage(ctx, userID, text)
        if err != nil {
            logger.Warn("Retrying to send message",
                zap.Int64("user_id", userID),
                zap.String("message", text),
                zap.Error(err))
            span.RecordError(err)
        }
        return err
    }, b)
}

func (c *Consumer) Stop() {
    close(c.stopChan)
    c.reader.Close()
    logger.Info("Consumer stopped")
}