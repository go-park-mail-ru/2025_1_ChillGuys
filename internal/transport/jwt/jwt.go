package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

// Секретный ключ для подписания токенов
var (
	SECRET = []byte("secret-key")
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
	secretKey []byte
}

// NewTokenator создает новый экземпляр Tokenator
func NewTokenator(secretKey []byte) *Tokenator {
	return &Tokenator{secretKey: secretKey}
}

// CreateJWT генерирует JWT токен для заданного userID и version
func (t *Tokenator) CreateJWT(userID string, version int) (string, error) {
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
	tokenString, err := token.SignedString(t.secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT Парсит и валидирует JWT-токен
func (t *Tokenator) ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return t.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Преобразуем claims в структуру Claims
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
