package app

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	http2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/http"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/search"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/suggestions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	addressrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	basketrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/basket"
	categoryrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/category"
	orderrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	productrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	searchrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/search"
	suggestionrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/suggestions"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	baskett "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	categoryt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	producttr "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	basketuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/basket"
	categoryuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/category"
	orderus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	searchus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/search"
	suggestionsus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/suggestions"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	usert "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user/http"
	csatt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/csat/http"
)

// App объединяет в себе все компоненты приложения.
type App struct {
	conf   *config.Config
	logger *logrus.Logger
	db     *sql.DB
	router *mux.Router
}

func OptionsRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
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
	minioClient, err := minio.NewMinioProvider(conf.MinioConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("minio initialization error: %w", err)
	}

	// Инициализация микросервисов
	authConn, err := grpc.Dial(
		"auth-service:50051",
		//":8010",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// Подключение к Redis
	redisSearchClient, err := redis.NewClient(conf.SearchRedisConfig)
	if err != nil {
		log.Fatalf("redis auth connection error: %v", err)
	}

	// Создаем Redis репозиторий
	redisSearchRepo := redis.NewSuggestionsRepository(redisSearchClient)

	authClient := auth.NewAuthServiceClient(authConn)

	userConn, err := grpc.Dial(
		"user-service:50052", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service connection error: %w", err)
	}
	userClient := user.NewUserServiceClient(userConn)

	userHandler := usert.NewUserHandler(userClient, conf)

	csatConn, err := grpc.Dial(
		"csat-service:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("csat service connection error: %w", err)
	}
	csatClient := csat.NewSurveyServiceClient(csatConn)

	csatHandler := csatt.NewCsatHandler(csatClient)

	// Инициализация репозиториев и use-case-ов.
	tokenator := jwt.NewTokenator(conf.JWTConfig)
	authHandler := http2.NewAuthHandler(authClient, conf)

	addressRepo := addressrepo.NewAddressRepository(db)
	addressUsecase := addressus.NewAddressUsecase(addressRepo)
	addressService := address.NewAddressHandler(addressUsecase, conf.GeoapifyConfig.APIKey)

	productRepo := productrepo.NewProductRepository(db)
	productUsecase := product.NewProductUsecase(productRepo)
	ProductService := producttr.NewProductService(productUsecase, minioClient)

	orderRepo := orderrepo.NewOrderRepository(db)
	orderUsecase := orderus.NewOrderUsecase(orderRepo)
	orderService := order.NewOrderService(orderUsecase)

	basketRepo := basketrepo.NewBasketRepository(db)
	basketUsecase := basketuc.NewBasketUsecase(basketRepo)
	basketService := baskett.NewBasketService(basketUsecase)

	categoryRepo := categoryrepo.NewCategoryRepository(db)
	categoryUsecase := categoryuc.NewCategoryUsecase(categoryRepo)
	categoryService := categoryt.NewCategoryService(categoryUsecase)

	suggestionsRepo := suggestionrepo.NewSuggestionsRepository(db)
	suggestionsUsecase := suggestionsus.NewSuggestionsUsecase(suggestionsRepo, redisSearchRepo)
	suggestionsService := suggestions.NewSuggestionsService(suggestionsUsecase)

	searchRepo := searchrepo.NewSearchRepository(db)
	searchUsecase := searchus.NewSearchUsecase(searchRepo)
	searchService := search.NewSearchService(searchUsecase, suggestionsUsecase)

	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter = apiRouter.PathPrefix("/v1").Subrouter()

	apiRouter.PathPrefix("/").HandlerFunc(OptionsRequest).Methods(http.MethodOptions)

	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.CORSMiddleware(next, conf.ServerConfig)
	})
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.LogRequest(logger, next)
	})

	// Маршруты для продуктов.
	productsRouter := apiRouter.PathPrefix("/products").Subrouter()
	{
		productsRouter.Handle("/batch",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(ProductService.GetProductsByIDs)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		productsRouter.HandleFunc("", ProductService.GetAllProducts).Methods(http.MethodGet)
		productsRouter.HandleFunc("/{id}", ProductService.GetProductByID).Methods(http.MethodGet)
		productsRouter.HandleFunc("/category/{id}", ProductService.GetProductsByCategory).Methods(http.MethodGet)
	}

	// Маршруты для категорий.
	catalogRouter := apiRouter.PathPrefix("/categories").Subrouter()
	{
		catalogRouter.HandleFunc("", categoryService.GetAllCategories).Methods(http.MethodGet)
	}

	suggestionsRouter := apiRouter.PathPrefix("/suggestions").Subrouter()
	{
		suggestionsRouter.HandleFunc("", suggestionsService.GetSuggestions).Methods(http.MethodGet)
	}

	searchRouter := apiRouter.PathPrefix("/search").Subrouter()
	{
		searchRouter.HandleFunc("", searchService.Search).Methods(http.MethodGet)
	}

	basketRouter := apiRouter.PathPrefix("/basket").Subrouter()
	{
		basketRouter.Handle("", middleware.JWTMiddleware(
			authClient,
			tokenator,
			http.HandlerFunc(basketService.Get)),
		).Methods(http.MethodGet)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(basketService.Add)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(basketService.Delete)),
				conf.CSRFConfig,
			)).Methods(http.MethodDelete)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(basketService.UpdateQuantity)),
				conf.CSRFConfig,
			)).Methods(http.MethodPatch)

		basketRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(basketService.Clear)),
				conf.CSRFConfig,
			)).Methods(http.MethodDelete)
	}

	productCoverRouter := apiRouter.PathPrefix("/cover").Subrouter()
	{
		productCoverRouter.HandleFunc("/upload", ProductService.CreateOne).Methods(http.MethodPost)
	}

	// Маршруты для аутентификации.
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
		authRouter.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
		authRouter.Handle("/logout",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(authHandler.Logout)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	// Маршруты для работы с пользователями.
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(userHandler.GetMe))).
			Methods(http.MethodGet)
		userRouter.Handle("/upload-avatar",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(userHandler.UploadAvatar)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-profile",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(userHandler.UpdateUserProfile)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-email",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(userHandler.UpdateUserEmail)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-password",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(userHandler.UpdateUserPassword)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	orderRouter := apiRouter.PathPrefix("/orders").Subrouter()
	{
		orderRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(orderService.CreateOrder)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		orderRouter.Handle("", middleware.JWTMiddleware(
			authClient,
			tokenator,
			http.HandlerFunc(orderService.GetOrders),
		)).Methods(http.MethodGet)
	}

	addressRouter := apiRouter.PathPrefix("/addresses").Subrouter()
	{
		addressRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(addressService.CreateAddress)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		addressRouter.Handle("", middleware.JWTMiddleware(
			authClient,
			tokenator,
			http.HandlerFunc(addressService.GetAddress),
		)).Methods(http.MethodGet)
		addressRouter.HandleFunc("/pickup-points", addressService.GetPickupPoints).Methods(http.MethodGet)
	}

	csatRouter := apiRouter.PathPrefix("").Subrouter()
	{
		csatRouter.Handle("csat/{name}", 
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.GetSurvey)),
				conf.CSRFConfig,
		)).Methods(http.MethodGet)

		csatRouter.Handle("csat",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.SubmitAnswer)),
				conf.CSRFConfig,
		)).Methods(http.MethodPost)

		csatRouter.Handle("/survey", 
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.GetAllSurveys)),
				conf.CSRFConfig,
		)).Methods(http.MethodGet)
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
func (a *App) Run() {

	server := &http.Server{
		Handler:      a.router,
		Addr:         fmt.Sprintf(":%s", a.conf.ServerConfig.Port),
		WriteTimeout: a.conf.ServerConfig.WriteTimeout,
		ReadTimeout:  a.conf.ServerConfig.ReadTimeout,
		IdleTimeout:  a.conf.ServerConfig.IdleTimeout,
	}

	a.logger.Infof("starting server on port %s", a.conf.ServerConfig.Port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Fatalf("server failed: %v", err)
	}
}
