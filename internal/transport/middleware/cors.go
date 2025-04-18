package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
)

func CORSMiddleware(next http.Handler, conf *config.ServerConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", conf.AllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AllowHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", conf.AllowCredentials)
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowOrigin)
		
		origin := r.Header.Get("Origin")

		fmt.Println(origin + " cors")
		if origin == "" {
			fmt.Println("return if origin")
			next.ServeHTTP(w, r)
			return
		}

		allowedOrigin := conf.AllowOrigin
		//allowedOrigin := "*"

		if allowedOrigin != "*" && allowedOrigin != origin {
			fmt.Println(allowedOrigin + "!=")
			http.Error(w, "CORS Origin not allowed", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
