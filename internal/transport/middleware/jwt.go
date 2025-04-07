package middleware

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// JWTMiddleware проверяет наличие и валидность JWT-токена в куках
func JWTMiddleware(tokenator *jwt.Tokenator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Получаем request_id из контекста
		reqID := ctx.Value(domains.ReqIDKey)
		if reqID == nil {
			reqID = "unknown"
		}

		// Логирование с request_id и remote_addr
		requestLogger := logrus.WithFields(logrus.Fields{
			"request_id":  reqID,
			"remote_addr": r.RemoteAddr,
			"path":        r.URL.Path,
		})

		cookieValue, err := r.Cookie(domains.Token.String())
		if err != nil {
			requestLogger.Warn("Missing or invalid token cookie")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Получаем токен из куки
		tokenString := cookieValue.Value

		// Если токен пустой, возвращаем ошибку
		if tokenString == "" {
			requestLogger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"ip":     r.RemoteAddr,
			}).Warn("Empty token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Вызываем ParseJWT через экземпляр Tokenator
		claims, err := tokenator.ParseJWT(tokenString)
		if err != nil {
			requestLogger.WithFields(logrus.Fields{
				"method": r.Method,
				"path":   r.URL.Path,
				"ip":     r.RemoteAddr,
				"error":  err.Error(),
			}).Error("Invalid token")

			response.SendJSONError(ctx, w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Проверяем, не истёк ли токен
		if claims.ExpiresAt < time.Now().Unix() {
			requestLogger.Warn("Token expired")
			response.SendJSONError(ctx, w, http.StatusUnauthorized, "Token expired")
			return
		}

		if !tokenator.VC.CheckUserVersion(ctx, claims.UserID, claims.Version) {
			requestLogger.Warn("Token is invalid or expired")
			response.SendJSONError(ctx, w, http.StatusUnauthorized, "Token is invalid or expired")
			return
		}

		// Передаем userID в контекст
		ctx = context.WithValue(ctx, domains.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
