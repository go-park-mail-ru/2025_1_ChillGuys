package grpc

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserIDInterceptor для обычных (unary) gRPC вызовов
func UserIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := extractUserID(ctx)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// UserIDStreamInterceptor для потоковых (stream) gRPC вызовов
func UserIDStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newCtx, err := extractUserID(ss.Context())
		if err != nil {
			return err
		}
		
		wrappedStream := &wrappedStream{
			ServerStream: ss,
			ctx:         newCtx,
		}
		return handler(srv, wrappedStream)
	}
}

// Обертка для ServerStream с измененным контекстом
type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

// extractUserID извлекает userID из метаданных и добавляет в контекст
func extractUserID(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, nil
	}

	userIDs := md.Get("user-id")
	if len(userIDs) == 0 {
		return ctx, nil
	}

	userID, err := uuid.Parse(userIDs[0])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user-id format")
	}

	if roles := md.Get("role"); len(roles) > 0 {
        ctx = context.WithValue(ctx, domains.RoleKey{}, roles[0])
    }

	return context.WithValue(ctx, domains.UserIDKey{}, userID.String()), nil
}