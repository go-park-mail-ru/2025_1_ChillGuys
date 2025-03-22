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
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	usecase2 "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		logger.WithFields(logrus.Fields{"error": "SERVER_PORT is not set"}).Error("SERVER_PORT is not set in the .env file")
		return
	}

	str, err := utils.GetConnectionString()
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

	userRepo := repository.NewUserRepository(db, logger)
	tokenator := jwt.NewTokenator(userRepo)
	userUsecase := usecase2.NewAuthUsecase(userRepo, tokenator, logger)
	userHandler := transport.NewAuthHandler(userUsecase, logger)

	productRepo := repository.NewProductRepository()
	productHandler := transport.NewProductHandler(productRepo, logger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(middleware.CORSMiddleware)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	productsRouter := router.PathPrefix("/products").Subrouter()
	{
		productsRouter.HandleFunc("/", productHandler.GetAllProducts).Methods("GET")
		productsRouter.HandleFunc("/{id}", productHandler.GetProductByID).Methods("GET")
		productsRouter.HandleFunc("/{id}/cover", productHandler.GetProductCover).Methods("GET")
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
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Infof("starting server on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Errorf("server error: %v", err)
	}
}
