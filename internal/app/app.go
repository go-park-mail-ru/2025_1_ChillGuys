package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	addressrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	authrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/auth"
	basketrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/basket"
	categoryrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/category"
	orderrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	productrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	userrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth"
	baskett "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	categoryt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	producttr "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	authus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
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
	conf   *config.Config
	logger *logrus.Logger
	db     *sql.DB
	router *mux.Router
	// Дополнительно можно добавить другие компоненты, если потребуется.
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

	// Инициализация репозиториев и use-case-ов.
	authRepo := authrepo.NewAuthRepository(db, logger)
	tokenator := jwt.NewTokenator(authRepo, conf.JWTConfig)
	authUsecase := authus.NewAuthUsecase(authRepo, tokenator, logger)
	authHandler := auth.NewAuthHandler(authUsecase, logger, conf)

	userRepository := userrepo.NewUserRepository(db)
	userUsecase := userus.NewUserUsecase(userRepository, tokenator, logger, minioClient)
	userService := user.NewUserHandler(userUsecase, logger, minioClient, conf)

	addressRepo := addressrepo.NewAddressRepository(db, logger)
	addressUsecase := addressus.NewAddressUsecase(addressRepo, logger)
	addressService := address.NewAddressHandler(addressUsecase, logger)

	productRepo := productrepo.NewProductRepository(db)
	productUsecase := product.NewProductUsecase(productRepo)
	ProductService := producttr.NewProductService(productUsecase, minioClient)

	orderRepo := orderrepo.NewOrderRepository(db)
	orderUsecase := orderus.NewOrderUsecase(orderRepo, logger)
	orderService := order.NewOrderService(orderUsecase, logger)

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

	router.PathPrefix("/").HandlerFunc(OptionsRequest).Methods(http.MethodOptions)

	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.CORSMiddleware(next, conf.ServerConfig)
	})
	apiRouter.Use(func(next http.Handler) http.Handler {
		return middleware.LogRequest(logger, next)
	})

	// Маршруты для продуктов.
	productsRouter := apiRouter.PathPrefix("/products").Subrouter()
	{
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

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketService.Add)),
		).Methods(http.MethodPost)

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketService.Delete)),
		).Methods(http.MethodDelete)

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketService.UpdateQuantity)),
		).Methods(http.MethodPatch)

		basketRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketService.Clear)),
		).Methods(http.MethodDelete)
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
		authRouter.Handle("/logout", middleware.JWTMiddleware(tokenator, http.HandlerFunc(authHandler.Logout))).
			Methods(http.MethodPost)
	}

	// Маршруты для работы с пользователями.
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.GetMe))).
			Methods(http.MethodGet)
		userRouter.Handle("/upload-avatar", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UploadAvatar))).
			Methods(http.MethodPost)
		userRouter.Handle("/update-profile", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserProfile))).
			Methods(http.MethodPost)
		userRouter.Handle("/update-email", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserEmail))).
			Methods(http.MethodPost)
		userRouter.Handle("/update-password", middleware.JWTMiddleware(tokenator, http.HandlerFunc(userService.UpdateUserPassword))).
			Methods(http.MethodPost)
	}

	orderRouter := router.PathPrefix("/orders").Subrouter()
	{
		orderRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(orderService.CreateOrder),
		)).Methods(http.MethodPost)
		orderRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(orderService.GetOrders),
		)).Methods(http.MethodGet)
	}

	addressRouter := router.PathPrefix("/addresses").Subrouter()
	{
		addressRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(addressService.CreateAddress),
		)).Methods(http.MethodPost)
		addressRouter.Handle("/", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(addressService.GetAddress),
		)).Methods(http.MethodGet)
		addressRouter.HandleFunc("/pickup-points", addressService.GetPickupPoints).Methods(http.MethodGet)
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
