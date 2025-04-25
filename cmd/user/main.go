package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	userrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	user "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user/grpc"
	us "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	grpcmw "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/grpc"
)

func main() {
	// Конфигурация
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	logger := logrus.New()

	// Подключение к Minio
	minioClient, err := minio.NewMinioProvider(conf.MinioConfig, logger)
	if err != nil {
		log.Fatalf("minio connection error: %v", err)
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

	tokenator := jwt.NewTokenator(conf.JWTConfig)

	// Инициализация репозиториев
	userRepo := userrepo.NewUserRepository(db)

	// Инициализация usecase
	userUsecase := us.NewUserUsecase(userRepo, tokenator, minioClient)

	// Создаем gRPC хендлер
	handler := user.NewUserGRPCHandler(userUsecase, minioClient)

	// Создаём gRPC сервер
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcmw.UserIDStreamInterceptor(),
		),
	)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gen.RegisterUserServiceServer(grpcServer, handler)

	fmt.Println("Starting user service on port :50052")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}