package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"
)

// JWTMiddleware проверяет наличие и валидность JWT-токена
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Warn("Missing Authorization header")
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Ожидаемый формат: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logrus.Warn("Invalid Authorization header format")
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Разбираем токен
		tokenString := parts[1]

		tokenator := jwt.NewTokenator([]byte("secret-key"))

		// Вызываем ParseJWT через экземпляр Tokenator
		claims, err := tokenator.ParseJWT(tokenString)
		if err != nil {
			logrus.Errorf("Invalid token: %v", err)
			utils.SendErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Проверяем, не истёк ли токен
		if claims.ExpiresAt < time.Now().Unix() {
			logrus.Warn("Token expired")
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Token expired")
			return
		}

		// Можно передавать UserID в контексте запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
