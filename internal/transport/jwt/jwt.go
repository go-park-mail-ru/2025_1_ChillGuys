package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

// Claims структура для данных токена
type JWTClaims struct {
	UserID    string
	Version   int
	ExpiresAt int64
	jwt.StandardClaims
}

// Tokenator структура для создания и парсинга токенов
type Tokenator struct {
}

// NewTokenator создает новый экземпляр Tokenator
func NewTokenator() *Tokenator {
	return &Tokenator{}
}

// CreateJWT генерирует JWT токен для заданного userID и version
func (t *Tokenator) CreateJWT(userID string, version int) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("TOKEN_SECRET")

	now := time.Now()
	expiration := now.Add(time.Hour * 24)

	claims := JWTClaims{
		UserID:    userID,
		Version:   version,
		ExpiresAt: expiration.Unix(),
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: expiration.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT парсит и валидирует JWT-токен
func (t *Tokenator) ParseJWT(tokenString string) (*JWTClaims, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secretKey := os.Getenv("TOKEN_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
