package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/app"
)

func main() {
	// Инициализация конфигурации.
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Создание экземпляра приложения.
	application, err := app.NewApp(conf)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	// Запуск приложения.
	if err := application.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
