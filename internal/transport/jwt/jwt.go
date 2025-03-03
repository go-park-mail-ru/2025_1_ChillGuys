package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	SECRET = []byte("secret-key")
)

type Claims struct {
	UserID    string
	Version   int
	ExpiresAt int64
	jwt.StandardClaims
}

func CreateJWT(userID string, version int) (string, error) {
	now := time.Now()
	expiration := now.Add(time.Hour * 24)

	claims := Claims{
		UserID:    userID,
		Version:   version,
		ExpiresAt: expiration.Unix(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: expiration.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	str, err := token.SignedString(SECRET)
	if err != nil {
		return "", err
	}
	return str, nil
}

// ParseJWT Парсит и валидирует JWT-токен
func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SECRET, nil
	})

	if err != nil {
		return nil, err
	}

	// Преобразуем claims в структуру Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
