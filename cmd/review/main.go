package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	reviewrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/review"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	review "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/review/grpc"
	cs "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/review"
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
	reviewRepo := reviewrepo.NewReviewRepository(db)

	// Инициализация usecase
	reviewUsecase := cs.NewReviewUsecase(reviewRepo)

	// Создаем gRPC хендлер
	handler := review.NewReviewGRPCServer(reviewUsecase)

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),
		),
	)

	lis, err := net.Listen("tcp", ":50054") 
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gen.RegisterReviewServiceServer(grpcServer, handler)

	fmt.Println("Starting review service on port :50054")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}