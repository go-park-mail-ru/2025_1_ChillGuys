package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

const (
	TokenLifeSpan = time.Hour * 24
)

type VersionChecker interface {
	CheckUserVersion(userID string, version int) bool
}

// JWTClaims структура для данных токена
type JWTClaims struct {
	UserID    string
	Version   int
	ExpiresAt int64
	jwt.StandardClaims
}

// Tokenator структура для создания и парсинга токенов
type Tokenator struct {
	sign string
	VC   VersionChecker
}

// NewTokenator создает новый экземпляр Tokenator
func NewTokenator(vc VersionChecker) *Tokenator {
	return &Tokenator{
		sign: os.Getenv("JWT_SIGNATURE"),
		VC:   vc,
	}
}

// CreateJWT генерирует JWT токен для заданного userID и version
func (t *Tokenator) CreateJWT(userID string, version int) (string, error) {
	now := time.Now()
	expiration := now.Add(TokenLifeSpan)

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

	return token.SignedString([]byte(t.sign))
}

// ParseJWT парсит и валидирует JWT-токен
func (t *Tokenator) ParseJWT(tokenString string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(t.sign), nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse jwt: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, models.ErrInvalidToken
}
