package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	csatrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/csat"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	csat "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/csat/grpc"
	cs "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/csat"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	grpcmw "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/grpc"
)

func main() {
	// Конфигурация
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Подключение к базе данных
	str, err := postgres.GetConnectionString(conf.DBConfig)
	if err != nil {
		log.Fatalf("db connection string error: %v", err)
	}
	db, err := sql.Open("postgres", str)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer db.Close()

	// Настройка пула соединений
	config.ConfigureDB(db, conf.DBConfig)

	// Инициализация репозиториев
	csatRepo := csatrepo.NewSurveyRepository(db)

	// Инициализация usecase
	csatUsecase := cs.NewCsatUsecase(csatRepo)

	// Создаем gRPC хендлер
	handler := csat.NewCsatGRPCHandler(csatUsecase)

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcmw.UserIDStreamInterceptor(),
		),
	)

	lis, err := net.Listen("tcp", ":50053") // Используем другой порт для CSAT сервиса
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gen.RegisterSurveyServiceServer(grpcServer, handler)

	fmt.Println("Starting CSAT service on port :50053")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}