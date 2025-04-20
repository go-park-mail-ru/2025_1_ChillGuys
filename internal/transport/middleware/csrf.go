package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
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

		//FIXME: включить для прода
		//jwtCookie, err := r.Cookie(string(domains.TokenCookieName))
		//if err != nil {
		//	http.Error(w, "JWT token required", http.StatusForbidden)
		//	return
		//}
		//jwtToken := jwtCookie.Value
		//
		//claims, err := tokenator.ParseJWT(jwtToken)
		//if err != nil {
		//	http.Error(w, "Invalid JWT token", http.StatusForbidden)
		//	return
		//}
		//
		//userID, err := uuid.Parse(claims.UserID)
		//if err != nil {
		//	http.Error(w, "Invalid user ID in token", http.StatusForbidden)
		//	return
		//}
		//
		//token := r.Header.Get(CSRFTokenHeader)
		//if token == "" {
		//	http.Error(w, "CSRF token missing", http.StatusForbidden)
		//	return
		//}
		//
		//valid, err := CheckCSRFToken(jwtToken, userID, token, cfg.SecretKey)
		//if err != nil || !valid {
		//	http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}

func CheckCSRFToken(tokenJWT string, userID uuid.UUID, tokenCSRF string, secretKey string) (bool, error) {
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
	data := fmt.Sprintf("%s:%s:%d", tokenJWT, userID.String(), csrfExp)
	h.Write([]byte(data))
	expectedMAC := h.Sum(nil)
	messageMAC, err := hex.DecodeString(csrfData[0])
	if err != nil {
		return false, fmt.Errorf("cannot hex decode tokenJWT")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}

func GenerateCSRFToken(tokenJWT string, userID uuid.UUID, secretKey string, tokenExpiry time.Duration) (string, error) {
	expiry := time.Now().Add(tokenExpiry).Unix()

	h := hmac.New(sha256.New, []byte(secretKey))
	data := fmt.Sprintf("%s:%s:%d", tokenJWT, userID.String(), expiry)
	h.Write([]byte(data))
	mac := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("%s:%d", mac, expiry), nil
}
