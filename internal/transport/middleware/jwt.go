package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

// JWTMiddleware проверяет наличие и валидность JWT-токена в куках
func JWTMiddleware(tokenator *jwt.Tokenator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieValue, err := r.Cookie(string(cookie.Token))
		if err != nil {
			logrus.Warn("Missing or invalid token cookie")
			response.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Получаем токен из куки
		tokenString := cookieValue.Value

		// Если токен пустой, возвращаем ошибку
		if tokenString == "" {
			logrus.Warn("Empty token")
			response.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Вызываем ParseJWT через экземпляр Tokenator
		claims, err := tokenator.ParseJWT(tokenString)
		if err != nil {
			logrus.Errorf("Invalid token: %v", err)
			response.SendErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		// Проверяем, не истёк ли токен
		if claims.ExpiresAt < time.Now().Unix() {
			logrus.Warn("Token expired")
			response.SendErrorResponse(w, http.StatusUnauthorized, "Token expired")
			return
		}

		ctx := r.Context()

		if !tokenator.VC.CheckUserVersion(ctx, claims.UserID, claims.Version) {
			response.SendErrorResponse(w, http.StatusUnauthorized, "Token is invalid or expired")
			return
		}

		ctx = context.WithValue(ctx, cookie.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
