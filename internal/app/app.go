package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres"
	basketrepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/basket"
	basketuc "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/basket"
	baskett "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/basket"
	product2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/product"
	user2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	product3 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"
	usecase2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() error {
	logger := logrus.New()

	// Получение конфигурации
	conf, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// Подключение базы данных
	str, err := postgres.GetConnectionString(conf.DBConfig)
	if err != nil {
		return fmt.Errorf("connection string error: %w", err)
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return fmt.Errorf("database connection error: %w", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Инициализация соединения с Minio
	minioClient, err := minio.NewMinioClient(conf.MinioConfig)
	if err != nil {
		return fmt.Errorf("minio initialization error: %w", err)
	}

	userRepo := user2.NewUserRepository(db, logger)
	tokenator := jwt.NewTokenator(userRepo, conf.JWTConfig)
	userUsecase := usecase2.NewAuthUsecase(userRepo, tokenator, logger, minioClient)
	userHandler := user.NewAuthHandler(userUsecase, logger, minioClient)

	productRepo := product2.NewProductRepository(db, logger)
	productUsecase := product.NewProductUsecase(logger, productRepo)
	productHandler := product3.NewProductHandler(productUsecase, logger, minioClient)

	basketRepo := basketrepo.NewBasketRepository(db, logger)
	basketUsecase := basketuc.NewBasketUsecase(logger, basketRepo)
	basketHandler := baskett.NewBasketService(basketUsecase, logger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router = router.PathPrefix("/v1").Subrouter()
	router.Use(func(next http.Handler) http.Handler {
		return middleware.CORSMiddleware(next, conf.ServerConfig)
	})
	router.Use(middleware.NewLoggerMiddleware(logger).LogRequest)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	productsRouter := router.PathPrefix("/products").Subrouter()
	{
		productsRouter.HandleFunc("/", productHandler.GetAllProducts).Methods("GET")
		productsRouter.HandleFunc("/{id}", productHandler.GetProductByID).Methods("GET")
		productsRouter.HandleFunc("/{id}/cover", productHandler.GetProductCover).Methods("GET")
		productsRouter.HandleFunc("/category/{id}", productHandler.GetProductsByCategory).Methods("GET")
	}

	catalogRouter := router.PathPrefix("/categories").Subrouter()
	{
		catalogRouter.HandleFunc("/", productHandler.GetAllCategories).Methods("GET")
	}

	basketRouter := router.PathPrefix("/basket").Subrouter()
	{
		basketRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketHandler.GetBasket)),
		).Methods(http.MethodGet)

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketHandler.AddProduct)),
		).Methods(http.MethodPost)

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketHandler.DeleteProduct)),
		).Methods(http.MethodDelete)

		basketRouter.Handle("/{id}", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketHandler.UpdateQuantity)),
		).Methods(http.MethodPut)

		basketRouter.Handle("", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(basketHandler.ClearBasket)),
		).Methods(http.MethodDelete)
	}

	productCoverRouter := router.PathPrefix("/cover").Subrouter()
	{
		productCoverRouter.HandleFunc("/upload", productHandler.CreateOne).Methods("POST")
		productsRouter.HandleFunc("/files/{objectID}", productHandler.GetOne).Methods("GET")
	}

	authRouter := router.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
		authRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
		authRouter.Handle("/logout", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(userHandler.Logout)),
		).Methods("POST")
	}

	userRouter := router.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(userHandler.GetMe)),
		).Methods("GET")

		userRouter.Handle("/upload-avatar", middleware.JWTMiddleware(
			tokenator,
			http.HandlerFunc(userHandler.UploadAvatar),
		)).Methods("POST")
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", conf.ServerConfig.Port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Infof("starting server on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
