package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	userrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	grpcmw "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/grpc"
	user "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user/grpc"
	us "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/gorilla/mux"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	// Инициализация токенатора
	tokenator := jwt.NewTokenator(conf.JWTConfig)

	// Инициализация репозиториев
	userRepo := userrepo.NewUserRepository(db)

	// Инициализация usecase
	userUsecase := us.NewUserUsecase(userRepo, tokenator, minioClient)

	// Создание gRPC хендлера
	handler := user.NewUserGRPCHandler(userUsecase, minioClient)

	// Создание middleware для метрик
	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceUserName)

	// Создаём gRPC сервер с цепочкой интерсепторов
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),             // ваш кастомный интерсептор
			metricsMw.ServerMetricsInterceptor,     // интерсептор метрик
			grpc_prometheus.UnaryServerInterceptor, // стандартный интерсептор prometheus
		),
		grpc.ChainStreamInterceptor(
			grpcmw.UserIDStreamInterceptor(),        // stream интерсептор
			grpc_prometheus.StreamServerInterceptor, // stream интерсептор prometheus
		),
	)

	// Регистрация метрик prometheus
	grpc_prometheus.Register(grpcServer)

	// Регистрируем сервис
	gen.RegisterUserServiceServer(grpcServer, handler)

	// Поднимаем HTTP-сервер для метрик
	go func() {
		router := mux.NewRouter()
		apiRouter := router.PathPrefix("/api").Subrouter()
		apiRouter = apiRouter.PathPrefix("/v1").Subrouter()
		apiRouter.PathPrefix("/metrics").Handler(promhttp.Handler())

		httpSrv := &http.Server{
			Addr:    ":8086",
			Handler: apiRouter,
		}

		log.Println("Metrics server starting on :8086")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("gRPC server starting on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
