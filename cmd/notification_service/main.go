//Точка входа микросервиса
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SeiFlow-3P2/notification_service/internal/app"
)

func main() {
	// Инициализация приложения
	a, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Запуск приложения
	if err := a.Run(); err != nil {
		log.Fatalf("Failed to run app: %v", err)
	}

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Received shutdown signal")
	a.Shutdown()
}