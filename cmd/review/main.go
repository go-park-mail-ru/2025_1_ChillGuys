package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	reviewrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/review"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	grpcmw "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/grpc"
	review "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/review/grpc"
	cs "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/review"
	"github.com/gorilla/mux"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
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

	// Инициализация middleware для метрик
	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceReviewName)

	// Создаём gRPC сервер с цепочкой интерсепторов
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmw.UserIDInterceptor(),             // ваш кастомный интерсептор
			metricsMw.ServerMetricsInterceptor,     // интерсептор метрик
			grpc_prometheus.UnaryServerInterceptor, // стандартный интерсептор prometheus
		),
	)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Регистрируем сервис
	gen.RegisterReviewServiceServer(grpcServer, handler)

	// Поднимаем HTTP-сервер для метрик
	go func() {
		router := mux.NewRouter()
		apiRouter := router.PathPrefix("/api").Subrouter()
		apiRouter = apiRouter.PathPrefix("/v1").Subrouter()
		apiRouter.PathPrefix("/metrics").Handler(promhttp.Handler())

		httpSrv := &http.Server{
			Addr:    ":8087",
			Handler: apiRouter,
		}

		log.Println("Metrics server starting on :8087")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server error: %v", err)
		}
	}()

	fmt.Println("Starting review service on port :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
