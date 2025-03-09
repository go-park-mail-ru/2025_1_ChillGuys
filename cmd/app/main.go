// @title ChillGuys API
// @version 1.0
// @description API for ChillGuys marketplace
// @host localhost:8080
// @BasePath /api

package main

import (
	"fmt"
	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	logger := logrus.New()

	userRepo := repository.NewUserRepository()
	tokenator := jwt.NewTokenator()
	userHandler := transport.NewAuthHandler(userRepo, logger, tokenator)

	productRepo := repository.NewProductRepo()
	productHandler := transport.NewProductHandler(productRepo, logger)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		logger.WithFields(logrus.Fields{"error": "SERVER_PORT is not set"}).Error("SERVER_PORT is not set in the .env file")
		return
	}

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
			http.HandlerFunc(userHandler.Logout)),
		).Methods("POST")
	}

	userRouter := router.PathPrefix("/users").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(
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
