//	@title			ChillGuys API
//	@version		1.0
//	@description	API for ChillGuys marketplace
//	@host			90.156.217.63:8081
//	@BasePath		/api

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	TokenAuth
//	@in							cookie
//	@name						token

package main

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	usecase2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase"
	"log"
	"net/http"
	"time"

	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
)

func main() {
	logger := logrus.New()

	// Получение конфигурации
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// Подключение базы данных
	str, err := utils.GetConnectionString(conf.DBConfig)
	if err != nil {
		logger.Error(err)
		return
	}
	db, err := sql.Open("postgres", str)
	if err != nil {
		logger.Error(err)
		return
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Инициализация соединения с Minio
	minioClient, err := minio.NewMinioClient(conf.MinioConfig)
	if err != nil {
		log.Fatalf("Minio initialization error: %v", err)
	}

	userRepo := repository.NewUserRepository(db, logger)
	tokenator := jwt.NewTokenator(userRepo, conf.JWTConfig)
	userUsecase := usecase2.NewAuthUsecase(userRepo, tokenator, logger, minioClient)
	userHandler := transport.NewAuthHandler(userUsecase, logger, minioClient)

	productRepo := repository.NewProductRepository(db, logger)
	productUsecase := usecase2.NewProductUsecase(logger, productRepo)
	productHandler := transport.NewProductHandler(productUsecase, logger, minioClient)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
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
		productsRouter.HandleFunc("/upload", productHandler.CreateOne).Methods("POST")
		router.HandleFunc("/files/{objectID}", productHandler.GetOne).Methods("GET")
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
		logger.Errorf("server error: %v", err)
	}
}
