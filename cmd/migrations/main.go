package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dsn, err := postgres.GetConnectionString(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	migrationsPath := "file://db/migrations"

	m, err := migrate.New(migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Error initializing migrations: %v", err)
	}

	if err = m.Up(); err != nil && err.Error() != "no change" {
		log.Fatalf("Error applying migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully.")
}
