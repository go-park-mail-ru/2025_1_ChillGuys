package middleware

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"net/http"
)

func CORSMiddleware(next http.Handler, conf *config.ServerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", conf.AllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AllowHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", conf.AllowCredentials)
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowOrigin)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
