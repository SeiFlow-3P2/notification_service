package logger

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init() error {
    // Настраиваем конфигурацию логгера
    config := zap.NewProductionConfig()
    // Логи в формате JSON
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.EncoderConfig.TimeKey = "timestamp"

    // Создаём логгер
    var err error
    logger, err = config.Build()
    if err != nil {
        return err
    }
    return nil
}

func Info(msg string, fields ...zap.Field) {
    logger.Info(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
    logger.Error(msg, append(fields, zap.Error(err))...)
}

func Warn(msg string, fields ...zap.Field) {
    logger.Warn(msg, fields...)
}

func Sync() {
    logger.Sync()
}