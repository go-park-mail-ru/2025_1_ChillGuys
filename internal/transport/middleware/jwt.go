package middleware

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// JWTMiddleware проверяет наличие и валидность JWT-токена в куках
func JWTMiddleware(tokenator *jwt.Tokenator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(string(utils.Token))
		if err != nil {
			logrus.Warn("Missing or invalid token cookie")
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Получаем токен из куки
		tokenString := cookie.Value

		// Если токен пустой, возвращаем ошибку
		if tokenString == "" {
			logrus.Warn("Empty token")
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

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

		if !tokenator.VC.CheckUserVersion(claims.UserID, claims.Version) {
			utils.SendErrorResponse(w, http.StatusUnauthorized, "Token is invalid or expired")
			return
		}

		// передаём UserID в контексте запроса
		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
