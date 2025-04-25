package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UserIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Получаем userID из метаданных gRPC
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if userIDs := md.Get("user-id"); len(userIDs) > 0 {
				if userID, err := uuid.Parse(userIDs[0]); err == nil {
					// Добавляем userID в контекст как раньше
					ctx = context.WithValue(ctx, domains.UserIDKey{}, userID.String())
				}
			}
		}
		return handler(ctx, req)
	}
}