package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewClient(broker string) (*kafka.Reader, error) {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     "calendar.notify",
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	}), nil
}