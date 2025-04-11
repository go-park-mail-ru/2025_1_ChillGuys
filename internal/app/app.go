package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	addressrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	orderrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	productrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	userrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	producttr "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	orderus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	userus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// App объединяет в себе все компоненты приложения.
type App struct {
	conf   *config.Config
	logger *logrus.Logger
	db     *sql.DB
	router *mux.Router
	// Дополнительно можно добавить другие компоненты, если потребуется.
}

// NewApp инициализирует приложение, создавая все необходимые компоненты.
func NewApp(conf *config.Config) (*App, error) {
	logger := logrus.New()

	// Подключение к базе данных.
	str, err := postgres.GetConnectionString(conf.DBConfig)
	if err != nil {
		return nil, fmt.Errorf("connection string error: %w", err)
	}
	db, err := sql.Open("postgres", str)
	if err != nil {
		return nil, fmt.Errorf("database connection error: %w", err)
	}

	// Применяем параметры пула соединений из конфигурации.
	config.ConfigureDB(db, conf.DBConfig)

	// Инициализация клиента Minio.
	ctx := context.Background()
	minioClient, err := minio.NewMinioClient(ctx, conf.MinioConfig)
	if err != nil {
		return nil, fmt.Errorf("minio initialization error: %w", err)
	}

	// Инициализация репозиториев и use-case-ов.
	userRepo := userrepo.NewUserRepository(db, logger)
	tokenator := jwt.NewTokenator(userRepo, conf.JWTConfig)
	userUsecase := userus.NewAuthUsecase(userRepo, tokenator, logger, minioClient)
	userHandler := user.NewAuthHandler(userUsecase, logger, minioClient, conf)

	addressRepo := addressrepo.NewAddressRepository(db, logger)
	addressUsecase := addressus.NewAddressUsecase(addressRepo, logger)
	addressHandler := address.NewAddressHandler(addressUsecase, logger)

	productRepo := productrepo.NewProductRepository(db, logger)
	productUsecase := product.NewProductUsecase(logger, productRepo)
	productHandler := producttr.NewProductHandler(productUsecase, logger, minioClient)

	orderRepo := orderrepo.NewOrderRepository(db, logger)
	orderUsecase := orderus.NewOrderUsecase(orderRepo, logger)
	orderHandler := order.NewOrderHandler(orderUsecase, logger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(func(next http.Handler) http.Handler {
		return middleware.CORSMiddleware(next, conf.ServerConfig)
	})
	router.Use(middleware.NewLoggerMiddleware(logger).LogRequest)

	// Подключение Swagger.
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Маршруты для продуктов.
	productsRouter := router.PathPrefix("/products").Subrouter()
	{
		productsRouter.HandleFunc("/", productHandler.GetAllProducts).Methods(http.MethodGet)
		productsRouter.HandleFunc("/{id}", productHandler.GetProductByID).Methods(http.MethodGet)
		productsRouter.HandleFunc("/{id}/cover", productHandler.GetProductCover).Methods(http.MethodGet)
		productsRouter.HandleFunc("/category/{id}", productHandler.GetProductsByCategory).Methods(http.MethodGet)
	}

	// Маршруты для категорий.
	catalogRouter := router.PathPrefix("/categories").Subrouter()
	{
		catalogRouter.HandleFunc("/", productHandler.GetAllCategories).Methods(http.MethodGet)
	}

	// Маршруты для загрузки обложек продукта.
	productCoverRouter := router.PathPrefix("/cover").Subrouter()
	{
		productCoverRouter.HandleFunc("/upload", productHandler.CreateOne).Methods(http.MethodPost)
		productsRouter.HandleFunc("/files/{objectID}", productHandler.GetOne).Methods(http.MethodGet)
	}

	// Маршруты для аутентификации.
	authRouter := router.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
		authRouter.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
		authRouter.Handle("/logout", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userHandler.Logout))).
			Methods(http.MethodPost)
	}

	// Маршруты для работы с пользователями.
	userRouter := router.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userHandler.GetMe))).
			Methods(http.MethodGet)
		userRouter.Handle("/upload-avatar", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userHandler.UploadAvatar))).
			Methods(http.MethodPost)
		userRouter.Handle("/update-profile", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userHandler.UpdateUserProfile))).
			Methods(http.MethodPost)
	}

	orderRouter := router.PathPrefix("/orders").Subrouter()
	{
		orderRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(orderHandler.CreateOrder),
		)).Methods(http.MethodPost)
		orderRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(orderHandler.GetOrders),
		)).Methods(http.MethodGet)
	}

	addressRouter := router.PathPrefix("/addresses").Subrouter()
	{
		addressRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(addressHandler.CreateAddress),
		)).Methods(http.MethodPost)
		addressRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(addressHandler.GetAddress),
		)).Methods(http.MethodGet)
		addressRouter.HandleFunc("/pickup-points", addressHandler.GetPickupPoints).Methods(http.MethodGet)
	}

	app := &App{
		conf:   conf,
		logger: logger,
		db:     db,
		router: router,
	}

	return app, nil
}

// Run запускает HTTP-сервер.
func (a *App) Run() error {
	server := &http.Server{
		Handler:      a.router,
		Addr:         fmt.Sprintf(":%s", a.conf.ServerConfig.Port),
		WriteTimeout: a.conf.ServerConfig.WriteTimeout,
		ReadTimeout:  a.conf.ServerConfig.ReadTimeout,
		IdleTimeout:  a.conf.ServerConfig.IdleTimeout,
	}

	a.logger.Infof("starting server on port %s", a.conf.ServerConfig.Port)
	return server.ListenAndServe()
}
