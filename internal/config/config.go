package config

import "os"

type Config struct {
	TelegramToken string
	KafkaBroker   string
	KafkaTopic    string
}

func Load() (*Config, error) {
	return &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		KafkaBroker:   os.Getenv("KAFKA_BROKER"),
		KafkaTopic:    "calendar.notify",
	}, nil
}