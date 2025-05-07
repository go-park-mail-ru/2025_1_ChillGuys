package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
)

const (
	CSRFTokenHeader = "X-CSRF-Token"
)

func CSRFMiddleware(tokenator *jwt.Tokenator, next http.Handler, cfg *config.CSRFConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		jwtCookie, err := r.Cookie(string(domains.TokenCookieName))
		if err != nil {
			http.Error(w, "JWT token required", http.StatusForbidden)
			return
		}
		jwtToken := jwtCookie.Value
		
		token := r.Header.Get(CSRFTokenHeader)
		if token == "" {
			http.Error(w, "CSRF token missing", http.StatusForbidden)
			return
		}
		
		valid, err := CheckCSRFToken(jwtToken, token, cfg.SecretKey)
		if err != nil || !valid {
			http.Error(w, "Invalid CSRF token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func CheckCSRFToken(tokenJWT string, tokenCSRF string, secretKey string) (bool, error) {
	csrfData := strings.Split(tokenCSRF, ":")
	if len(csrfData) != 2 {
		return false, fmt.Errorf("bad tokenJWT data")
	}

	csrfExp, err := strconv.ParseInt(csrfData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("bad tokenJWT time")
	}

	if csrfExp < time.Now().Unix() {
		return false, fmt.Errorf("tokenJWT expired")
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	data := fmt.Sprintf("%s:%d", tokenJWT, csrfExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(csrfData[0])
	if err != nil {
		return false, fmt.Errorf("cannot hex decode tokenJWT")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}

func GenerateCSRFToken(tokenJWT string, secretKey string, tokenExpiry time.Duration) (string, error) {
	expiry := time.Now().Add(tokenExpiry).Unix()

	h := hmac.New(sha256.New, []byte(secretKey))
	data := fmt.Sprintf("%s:%d", tokenJWT, expiry)
	h.Write([]byte(data))
	mac := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%d", mac, expiry), nil
}

