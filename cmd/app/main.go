package main

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	logger := logrus.New()

	userRepo := repository.NewUserRepository()
	tokenator := jwt.NewTokenator()
	userHandler := transport.NewAuthHandler(userRepo, logger, tokenator)

	productRepo := repository.NewProductRepo()
	productHandler := transport.NewProductHandler(productRepo)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(middleware.CORSMiddleware)

	productsRouter := router.PathPrefix("/products").Subrouter()
	{
		productsRouter.HandleFunc("/", productHandler.GetAllProducts).Methods("GET")
		productsRouter.HandleFunc("/{id}", productHandler.GetProductByID).Methods("GET")
	}

	authRouter := router.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/login", userHandler.Login).Methods("POST")
		authRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
		authRouter.Handle("/logout", middleware.JWTMiddleware(
			http.HandlerFunc(userHandler.Logout)),
		).Methods("POST")
	}

	userRouter := router.PathPrefix("/user").Subrouter()
	{
		userRouter.Handle("/me", middleware.JWTMiddleware(
			http.HandlerFunc(userHandler.GetMe)),
		).Methods("GET")
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Infof("starting server on port %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		logger.Errorf("server error: %v", err)
	}
}
