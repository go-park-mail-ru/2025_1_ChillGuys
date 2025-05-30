package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	http2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth/http"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/csat"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/recommendation"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/search"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/suggestions"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	addressrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	adminrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/admin"
	basketrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/basket"
	categoryrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/category"
	promorepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/promo"
	orderrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/order"
	productrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	recrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/recommendation"
	searchrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/search"
	sellerrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/seller"
	suggestionrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/suggestions"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/address"
	admint "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/admin"
	baskett "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	categoryt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/category"
	csatt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/csat/http"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/order"
	producttr "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	promot "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/promo"
	notificationt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/notification"
	notificationuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/notification"
	motificationrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/notification"
	reviewt "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/review/http"
	sellert "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/seller"
	usert "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user/http"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"
	adminuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/admin"
	promouc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/promo"
	basketuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/basket"
	categoryuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/category"
	orderus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/order"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	recus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/recommendation"
	searchus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/search"
	selleruc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/seller"
	suggestionsus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/suggestions"
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

	reviewConn, err := grpc.Dial(
		"review-service:50054",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("review service connection error: %w", err)
	}
	reviewClient := review.NewReviewServiceClient(reviewConn)

	reviewHandler := reviewt.NewReviewHandler(reviewClient)

	// Инициализация репозиториев и use-case-ов.
	tokenator := jwt.NewTokenator(conf.JWTConfig)
	authHandler := http2.NewAuthHandler(authClient, conf)

	addressRepo := addressrepo.NewAddressRepository(db)
	addressUsecase := addressus.NewAddressUsecase(addressRepo)
	addressService := address.NewAddressHandler(addressUsecase, conf.GeoapifyConfig.APIKey)

	productRepo := productrepo.NewProductRepository(db)
	productUsecase := product.NewProductUsecase(productRepo)
	ProductService := producttr.NewProductService(productUsecase, minioClient)

	basketRepo := basketrepo.NewBasketRepository(db)
	basketUsecase := basketuc.NewBasketUsecase(basketRepo)
	basketService := baskett.NewBasketService(basketUsecase)

	categoryRepo := categoryrepo.NewCategoryRepository(db)
	categoryUsecase := categoryuc.NewCategoryUsecase(categoryRepo)
	categoryService := categoryt.NewCategoryService(categoryUsecase)

	suggestionsRepo := suggestionrepo.NewSuggestionsRepository(db)
	suggestionsUsecase := suggestionsus.NewSuggestionsUsecase(suggestionsRepo, redisSearchRepo)
	suggestionsService := suggestions.NewSuggestionsService(suggestionsUsecase)

	adminRepo := adminrepo.NewAdminRepository(db)
	adminUsecase := adminuc.NewAdminUsecase(adminRepo, redisSearchRepo, productRepo)
	adminService := admint.NewAdminService(adminUsecase)

	sellerRepo := sellerrepo.NewSellerRepository(db)
	sellerUsecase := selleruc.NewSellerUsecase(sellerRepo)
	sellerService := sellert.NewSellerHandler(sellerUsecase, minioClient)

	searchRepo := searchrepo.NewSearchRepository(db)
	searchUsecase := searchus.NewSearchUsecase(searchRepo)
	searchService := search.NewSearchService(searchUsecase, suggestionsUsecase)


	promoRepo := promorepo.NewPromoRepository(db)
	promoUsecase := promouc.NewPromoUsecase(promoRepo)
	promoService := promot.NewPromoService(promoUsecase)

	notificationRepo := motificationrepo.NewNotificationRepository(db)
	notificationUsecase := notificationuc.NewNotificationUsecase(notificationRepo)
	notificationService := notificationt.NewNotificationService(notificationUsecase)

	orderRepo := orderrepo.NewOrderRepository(db)
	orderUsecase := orderus.NewOrderUsecase(orderRepo, promoRepo, notificationRepo)
	orderService := order.NewOrderService(orderUsecase)

	recommendationRepo := recrepo.NewRecommendationRepository(db)
	recommendationUsecase := recus.NewRecommendationUsecase(productUsecase, recommendationRepo)
	recommendationServise := recommendation.NewRecommendationService(recommendationUsecase)


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

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceMainName)
	apiRouter.PathPrefix("/metrics").Handler(promhttp.Handler())
	apiRouter.Use(metricsMw.LogMetrics)

	// Маршруты для продуктов.
	productsRouter := apiRouter.PathPrefix("").Subrouter()
	{
		productsRouter.Handle("/products/batch",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(ProductService.GetProductsByIDs)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
		productsRouter.HandleFunc("/products/{offset}", ProductService.GetAllProducts).Methods(http.MethodGet)
		productsRouter.HandleFunc("/product/{id}", ProductService.GetProductByID).Methods(http.MethodGet)
		productsRouter.HandleFunc("/products/category/{id}/{offset}", ProductService.GetProductsByCategory).Methods(http.MethodGet)

		productsRouter.Handle("/add",
			http.HandlerFunc(ProductService.AddProduct),
		).Methods(http.MethodPost)
	}

	// Маршруты для категорий.
	catalogRouter := apiRouter.PathPrefix("/categories").Subrouter()
	{
		catalogRouter.HandleFunc("", categoryService.GetAllCategories).Methods(http.MethodGet)
		catalogRouter.HandleFunc("/{id}", categoryService.GetAllSubcategories).Methods(http.MethodGet)
	}

	subcategoryRouter := apiRouter.PathPrefix("/subcategory").Subrouter()
	{
		subcategoryRouter.HandleFunc("/{id}", categoryService.GetNameSubcategory).Methods(http.MethodGet)
	}

	suggestionsRouter := apiRouter.PathPrefix("/suggestions").Subrouter()
	{
		suggestionsRouter.HandleFunc("", suggestionsService.GetSuggestions).Methods(http.MethodPost)
	}

	searchRouter := apiRouter.PathPrefix("/search").Subrouter()
	{
		searchRouter.HandleFunc("/sort/{offset}", searchService.SearchWithFilterAndSort).Methods(http.MethodPost)
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

		userRouter.Handle("/update-role",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("buyer")(
						http.HandlerFunc(userHandler.BecomeSeller),
					),
				),
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
		csatRouter.Handle("/csat/{name}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.GetSurvey)),
				conf.CSRFConfig,
			)).Methods(http.MethodGet)

		csatRouter.Handle("/csat",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.SubmitAnswer)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		csatRouter.Handle("/survey",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.GetAllSurveys)),
				conf.CSRFConfig,
			)).Methods(http.MethodGet)

		csatRouter.Handle("/stat/{surveyId}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(csatHandler.GetSurveyStatistics)),
				conf.CSRFConfig,
			)).Methods(http.MethodGet)
	}

	reviewRouter := apiRouter.PathPrefix("/review").Subrouter()
	{
		reviewRouter.HandleFunc("", reviewHandler.Get).Methods(http.MethodPost)

		reviewRouter.Handle("/add",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(reviewHandler.Add)),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	notificationRouter := apiRouter.PathPrefix("/notification").Subrouter()
	{
		notificationRouter.Handle("/count",
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(notificationService.GetUnreadCount)),
			).Methods(http.MethodGet)

		notificationRouter.Handle("/{offset}",
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(notificationService.GetUserNotifications)),
			).Methods(http.MethodGet)

		notificationRouter.Handle("/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator, http.HandlerFunc(notificationService.MarkAsRead)),
				conf.CSRFConfig,
			)).Methods(http.MethodPatch)
	}

	adminRouter := apiRouter.PathPrefix("/admin").Subrouter()
	{
		adminRouter.Handle("/products/{offset}",
			middleware.JWTMiddleware(authClient, tokenator,
				middleware.RoleMiddleware("admin")(
					http.HandlerFunc(adminService.GetPendingProducts),
				),
			),
		).Methods(http.MethodGet)

		adminRouter.Handle("/users/{offset}",
			middleware.JWTMiddleware(authClient, tokenator,
				middleware.RoleMiddleware("admin")(
					http.HandlerFunc(adminService.GetPendingUsers),
				),
			),
		).Methods(http.MethodGet)

		adminRouter.Handle("/product/update",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("admin")(
						http.HandlerFunc(adminService.UpdateProductStatus),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		adminRouter.Handle("/user/update",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("admin")(
						http.HandlerFunc(adminService.UpdateUserRole),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	promoRouter := apiRouter.PathPrefix("/promo").Subrouter()
	{
		promoRouter.Handle("/",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("admin")(
						http.HandlerFunc(promoService.Create),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		promoRouter.Handle("/{offset}",
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("admin")(
						http.HandlerFunc(promoService.GetAll),
					),
			)).Methods(http.MethodGet)

		promoRouter.Handle("/check",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
						http.HandlerFunc(promoService.CheckPromoCode),
					),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	warehouseRouter := apiRouter.PathPrefix("/warehouse").Subrouter()
	{
		warehouseRouter.Handle("/get",
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("warehouseman")(
						http.HandlerFunc(orderService.GetOrdersPlaced),
					),
			)).Methods(http.MethodGet)

		warehouseRouter.Handle("/update",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("warehouseman")(
						http.HandlerFunc(orderService.UpdateStatus),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)
	}

	sellerRouter := apiRouter.PathPrefix("/seller").Subrouter()
	{
		sellerRouter.Handle("/add-product",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("seller")(
						http.HandlerFunc(sellerService.AddProduct),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		sellerRouter.Handle("/add-image/{id}",
			middleware.CSRFMiddleware(tokenator,
				middleware.JWTMiddleware(authClient, tokenator,
					middleware.RoleMiddleware("seller")(
						http.HandlerFunc(sellerService.UploadProductImage),
					),
				),
				conf.CSRFConfig,
			)).Methods(http.MethodPost)

		sellerRouter.Handle("/products/{offset}",
			middleware.JWTMiddleware(authClient, tokenator,
				middleware.RoleMiddleware("seller")(
					http.HandlerFunc(sellerService.GetSellerProducts),
				),
			),
		).Methods(http.MethodGet)
	}

	recommendationRouter := apiRouter.PathPrefix("/recommendation").Subrouter()
	{
		recommendationRouter.HandleFunc("/{id}", recommendationServise.GetRecommendations).Methods(http.MethodGet)
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