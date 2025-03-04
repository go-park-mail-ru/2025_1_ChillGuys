package main

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	logger := logrus.New()

	userRepo := repository.NewUserRepository()
	authHandler := auth.NewAuthHandler(userRepo, logger)

	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Use(middleware.CORSMiddleware)

	authRouter := router.PathPrefix("/auth").Subrouter()
	{
		authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
		authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
		authRouter.Handle("/logout", middleware.JWTMiddleware(
			http.HandlerFunc(authHandler.Logout)),
		).Methods("POST")
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	logger.Infof("starting server on port %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
