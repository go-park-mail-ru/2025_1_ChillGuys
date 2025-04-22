package app

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	http2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/http"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
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
	userrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	baskett "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	categoryt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	producttr "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	basketuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/basket"
	categoryuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/category"
	orderus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	userus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// App объединяет в себе все компоненты приложения.
type App struct {
	conf        *config.Config
	logger      *logrus.Logger
	db          *sql.DB
	redisClient *redis.Client
	router      *mux.Router
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

	// Подключение к Redis
	redisClient, err := redis.NewClient(conf.RedisConfig)
	if err != nil {
		log.Fatalf("redis connection error: %v", err)
	}

	// Создаем Redis репозиторий
	redisAuthRepo := redis.NewAuthRepository(redisClient, conf.JWTConfig)

	// Применяем параметры пула соединений из конфигурации.
	config.ConfigureDB(db, conf.DBConfig)

	// Инициализация клиента Minio.
	minioClient, err := minio.NewMinioProvider(conf.MinioConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("minio initialization error: %w", err)
	}

	// Инициализация микросервисов
	authConn, err := grpc.Dial(
		"localhost:50051",
		//":8010",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	authClient := auth.NewAuthServiceClient(authConn)

	// Инициализация репозиториев и use-case-ов.
	tokenator := jwt.NewTokenator(redisAuthRepo, conf.JWTConfig)
	authHandler := http2.NewAuthHandler(authClient, conf)

	userRepository := userrepo.NewUserRepository(db)
	userUsecase := userus.NewUserUsecase(userRepository, tokenator, minioClient)
	userService := user.NewUserHandler(userUsecase, minioClient, conf)

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
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(ProductService.GetProductsByIDs)),
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

	basketRouter := apiRouter.PathPrefix("/basket").Subrouter()
	{
		basketRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketService.Get)),
		).Methods(http.MethodGet)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(basketService.Add)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(basketService.Delete)),
				conf.CSRFConfig,
			)).Methods(http.MethodDelete)

		basketRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(basketService.UpdateQuantity)),
				conf.CSRFConfig,
			)).Methods(http.MethodPatch)

		basketRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(basketService.Clear)),
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
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(authHandler.Logout)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	// Маршруты для работы с пользователями.
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.GetMe))).
			Methods(http.MethodGet)
		userRouter.Handle("/upload-avatar",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UploadAvatar)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-profile",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserProfile)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-email",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserEmail)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		userRouter.Handle("/update-password",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserPassword)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	orderRouter := apiRouter.PathPrefix("/orders").Subrouter()
	{
		orderRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(orderService.CreateOrder)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		orderRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(orderService.GetOrders),
		)).Methods(http.MethodGet)
	}

	addressRouter := apiRouter.PathPrefix("/addresses").Subrouter()
	{
		addressRouter.Handle("",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(tokenator, http.HandlerFunc(addressService.CreateAddress)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		addressRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(addressService.GetAddress),
		)).Methods(http.MethodGet)
		addressRouter.HandleFunc("/pickup-points", addressService.GetPickupPoints).Methods(http.MethodGet)
	}

	app := &App{
		conf:        conf,
		logger:      logger,
		db:          db,
		redisClient: redisClient,
		router:      router,
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
