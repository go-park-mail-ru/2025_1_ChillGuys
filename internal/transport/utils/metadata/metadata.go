package metadata

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"google.golang.org/grpc/metadata"
	"strings"
)

const jwtKey = "authorization"

// ExtractJWTFromContext извлекает JWT токен из gRPC контекста
// Формат ожидаемого заголовка: "Bearer <token>"
func ExtractJWTFromContext(ctx context.Context) (string, error) {
	// 1. Получаем метаданные из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errs.ErrNoMetadata
	}

	// 2. Извлекаем заголовок Authorization
	authHeaders := md.Get(jwtKey)
	if len(authHeaders) == 0 {
		return "", errs.ErrNoAuthHeader
	}

	// 3. Проверяем формат "Bearer <token>"
	tokenHeader := authHeaders[0]
	if len(tokenHeader) < 7 || !strings.HasPrefix(tokenHeader, "Bearer ") {
		return "", errs.ErrInvalidAuthFormat
	}

	// 4. Возвращаем чистый токен (без "Bearer ")
	return tokenHeader[7:], nil
}

func InjectJWTIntoContext(ctx context.Context, token string) context.Context {
	md := metadata.New(map[string]string{jwtKey: "Bearer " + token})
	return metadata.NewOutgoingContext(ctx, md)
}
