package app

import (
    "context"
    "log"
    "net"
    "net/http"
    "time"

    "notification_service/config"
    "notification_service/consumer"
    "notification_service/kafka"
    "notification_service/logger"
    "notification_service/metrics"
    "notification_service/telegram"
    "notification_service/telemetry"

    "github.com/segmentio/kafka-go"
    "google.golang.org/grpc"
)

type App struct {
    consumer   *consumer.Consumer
    grpcServer *grpc.Server
    httpServer *http.Server
}

func NewApp() (*App, error) {
    // Инициализация логгера
    if err := logger.Init(); err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()

    // Инициализация телеметрии
    if err := telemetry.Init(); err != nil {
        logger.Error("Failed to initialize telemetry", err)
        return nil, err
    }

    // Загрузка конфигурации
    cfg, err := config.Load()
    if err != nil {
        logger.Error("Failed to load configuration", err)
        return nil, err
    }

    // Инициализация Kafka
    kafkaClient, err := kafka.NewClient(cfg.KafkaBroker)
    if err != nil {
        logger.Error("Failed to initialize Kafka client", err)
        return nil, err
    }

    // Инициализация Telegram клиента
    tgClient, err := telegram.NewClient(cfg.TelegramToken)
    if err != nil {
        logger.Error("Failed to initialize Telegram client", err)
        return nil, err
    }

    // Инициализация потребителя
    consumer, err := consumer.NewConsumer(kafkaClient, tgClient, cfg.KafkaTopic)
    if err != nil {
        logger.Error("Failed to initialize consumer", err)
        return nil, err
    }

    // Настройка gRPC сервера
    grpcServer := grpc.NewServer()

    // Настройка HTTP сервера для Prometheus
    httpServer := &http.Server{
        Addr:    ":9091",
        Handler: metrics.NewHandler(),
    }

    return &App{
        consumer:   consumer,
        grpcServer: grpcServer,
        httpServer: httpServer,
    }, nil
}

func (a *App) Run() error {
    // Запуск потребителя в горутине
    go a.consumer.Start()

    // Запуск gRPC сервера
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        logger.Error("Failed to listen on gRPC port", err)
        return err
    }
    go a.grpcServer.Serve(lis)
    logger.Info("gRPC server started", zap.String("port", ":50051"))

    // Запуск HTTP сервера для Prometheus
    go func() {
        if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("HTTP server failed", err)
        }
    }()
    logger.Info("HTTP server for Prometheus started", zap.String("port", ":9091"))

    return nil
}

func (a *App) Shutdown() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    a.consumer.Stop()
    a.grpcServer.GracefulStop()
    a.httpServer.Shutdown(ctx)
    logger.Info("Application shut down gracefully")
}