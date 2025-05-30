package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
)

func CORSMiddleware(next http.Handler, conf *config.ServerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", conf.AllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AllowHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", conf.AllowCredentials)
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowOrigin)

		allowedOrigin := conf.AllowOrigin
		//allowedOrigin := "*"

		if allowedOrigin != "*" && allowedOrigin != origin {
			http.Error(w, "CORS Origin not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
