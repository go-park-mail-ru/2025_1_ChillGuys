package main

import (
	"database/sql"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	authrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	auth "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/grpc"
	auth2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	au "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Конфигурация
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Подключение к Redis для аутентификации
	redisAuthClient, err := redis.NewClient(conf.AuthRedisConfig)
	if err != nil {
		log.Fatalf("redis auth connection error: %v", err)
	}

	// Создаем Redis репозиторий для аутентификации
	redisAuthRepo := redis.NewAuthRepository(redisAuthClient, conf.JWTConfig)

	// Подключение к базе данных
	str, err := postgres.GetConnectionString(conf.DBConfig)
	if err != nil {
		log.Fatalf("db error: %v", err)
		return
	}
	db, err := sql.Open("postgres", str)
	if err != nil {
		log.Fatalf("db error: %v", err)
		return
	}
	defer db.Close()

	// Применяем параметры пула соединений из конфигурации
	config.ConfigureDB(db, conf.DBConfig)

	// Создание токенатора JWT
	tokenator := jwt.NewTokenator(conf.JWTConfig)

	// Инициализация репозиториев
	authRepo := authrepo.NewAuthRepository(db)

	// Инициализация usecase с Redis репозиторием
	authUsecase := au.NewAuthUsecase(authRepo, redisAuthRepo, tokenator)

	// Создаем хендлер с передачей всех необходимых зависимостей
	handler := auth.NewAuthGRPCHandler(authUsecase, redisAuthRepo, tokenator)

	// Создаём сервер
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Регистрируем сервис
	auth2.RegisterAuthServiceServer(grpcServer, handler)

	log.Println("gRPC server starting on :50051")

	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
