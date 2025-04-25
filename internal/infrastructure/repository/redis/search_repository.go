package redis

import (
	"context"
	"fmt"
)

const (
	CategoryNamesKey = "suggestions:categories"
	ProductNamesKey  = "suggestions:products"
)

type SuggestionsRepository struct {
	client *Client
}

func NewSuggestionsRepository(client *Client) *SuggestionsRepository {
	return &SuggestionsRepository{client: client}
}

// AddSuggestionsByKey добавляет список строк в Redis по указанному ключу.
func (r *SuggestionsRepository) AddSuggestionsByKey(ctx context.Context, key string, names []string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection error: %w", err)
	}

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete old suggestions: %w", err)
	}

	values := make([]interface{}, 0, len(names))
	for _, name := range names {
		values = append(values, name)
	}

	if err := r.client.SAdd(ctx, key, values...).Err(); err != nil {
		return fmt.Errorf("failed to add suggestions to Redis: %w", err)
	}

	return nil
}

// GetSuggestionsByKey получает список строк из Redis по указанному ключу.
func (r *SuggestionsRepository) GetSuggestionsByKey(ctx context.Context, key string) ([]string, error) {
	if r.client == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	names, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from Redis: %w", err)
	}

	return names, nil
}
