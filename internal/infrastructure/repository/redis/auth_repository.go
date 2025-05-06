package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
)

const (
	userTokensPrefix = "user_id:"
)

type AuthRepository struct {
	client *Client
	cfg    *config.JWTConfig
}

func NewAuthRepository(client *Client, jwtCfg *config.JWTConfig) *AuthRepository {
	return &AuthRepository{
		client: client,
		cfg:    jwtCfg,
	}
}

// AddToBlacklist добавляет токен в список недействительных токенов пользователя
func (r *AuthRepository) AddToBlacklist(ctx context.Context, userID, token string) error {
	expiration := time.Until(time.Now().Add(r.cfg.TokenLifeSpan))
	userKey := fmt.Sprintf("%s%s", userTokensPrefix, userID)

	// Добавляем токен в множество пользователя
	if err := r.client.SAdd(ctx, userKey, token).Err(); err != nil {
		return fmt.Errorf("failed to add token to user's blacklist: %w", err)
	}

	// Обновляем TTL для множества, чтобы после жизни самого старшего токена ключ очистился
	if err := r.client.Expire(ctx, userKey, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration for user's blacklist: %w", err)
	}

	return nil
}

// IsInBlacklist проверяет, находится ли токен в черном списке
func (r *AuthRepository) IsInBlacklist(ctx context.Context, userID, token string) (bool, error) {
	if r.client == nil {
		return false, fmt.Errorf("redis client is not initialized")
	}
	
	userKey := fmt.Sprintf("%s%s", userTokensPrefix, userID)
	isMember, err := r.client.SIsMember(ctx, userKey, token).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check token in blacklist: %w", err)
	}
	return isMember, nil
}
