package config

import (
    "github.com/joho/godotenv"
    "os"
    "log"
)

type Config struct {
    TelegramToken string
    KafkaBroker   string
    KafkaTopic    string
}

func Load() (*Config, error) {
    // Загружаем переменные из .env
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: No .env file found, using environment variables only: %v", err)
    }

    return &Config{
        TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
        KafkaBroker:   os.Getenv("KAFKA_BROKER"),
        KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
    }, nil
}