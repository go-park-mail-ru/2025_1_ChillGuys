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
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	grpcmw "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/grpc"
	au "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/gorilla/mux"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
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

	// Инициализация middleware
	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceAuthName)

	// Создаём сервер с цепочкой интерцепторов
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),             // 1. Интерцептор для работы с UserID
			metricsMw.ServerMetricsInterceptor,     // 2. Интерцептор для метрик
			grpc_prometheus.UnaryServerInterceptor, // 3. Стандартный интерцептор метрик
		),
	)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Регистрируем сервис
	auth2.RegisterAuthServiceServer(grpcServer, handler)

	// Поднимаем HTTP-сервер для метрик на другом порту
	go func() {
		router := mux.NewRouter()
		apiRouter := router.PathPrefix("/api").Subrouter()
		apiRouter = apiRouter.PathPrefix("/v1").Subrouter()
		apiRouter.PathPrefix("/metrics").Handler(promhttp.Handler())

		httpSrv := &http.Server{
			Addr:    ":8085",
			Handler: apiRouter,
		}

		log.Println("metrics server starting on :8085")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	log.Println("gRPC server starting on :50051")
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
