package redis

import (
	"context"
	"fmt"
)

const (
	CategoryNamesKey = "suggestions:categories"
	ProductNamesKey  = "suggestions:products"
)

func categoryKey(categoryID string) string {
	return fmt.Sprintf("suggestions:products:category:%s", categoryID)
}

type SuggestionsRepository struct {
	client *Client
}

func NewSuggestionsRepository(client *Client) *SuggestionsRepository {
	return &SuggestionsRepository{client: client}
}

// AddSuggestionsByKey добавляет список строк в Redis по указанному ключу, сохраняя существующие записи
func (r *SuggestionsRepository) AddSuggestionsByKey(ctx context.Context, key string, names []string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection error: %w", err)
	}

	// Получаем текущие значения
	current, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get current suggestions: %w", err)
	}

	// Создаем мапу для быстрой проверки существующих значений
	existing := make(map[string]struct{})
	for _, v := range current {
		existing[v] = struct{}{}
	}

	// Добавляем только новые значения
	var newValues []interface{}
	for _, name := range names {
		if _, found := existing[name]; !found {
			newValues = append(newValues, name)
		}
	}

	if len(newValues) > 0 {
		if err := r.client.SAdd(ctx, key, newValues...).Err(); err != nil {
			return fmt.Errorf("failed to add new suggestions to Redis: %w", err)
		}
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

// AddProductSuggestionsByCategory сохраняет продукты для конкретной категории.
func (r *SuggestionsRepository) AddProductSuggestionsByCategory(ctx context.Context, categoryID string, names []string) error {
	return r.AddSuggestionsByKey(ctx, categoryKey(categoryID), names)
}

// GetProductSuggestionsByCategory получает продукты для конкретной категории.
func (r *SuggestionsRepository) GetProductSuggestionsByCategory(ctx context.Context, categoryID string) ([]string, error) {
	return r.GetSuggestionsByKey(ctx, categoryKey(categoryID))
}

// RemoveSuggestion удаляет конкретное значение из набора
func (r *SuggestionsRepository) RemoveSuggestion(ctx context.Context, key string, value string) error {
	if r.client == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	if err := r.client.SRem(ctx, key, value).Err(); err != nil {
		return fmt.Errorf("failed to remove suggestion: %w", err)
	}

	return nil
}

// GetSuggestionsByKeyPaginated получает список строк из Redis с пагинацией
// pageNum - номер страницы (начиная с 0)
// limit - количество элементов на странице
func (r *SuggestionsRepository) GetSuggestionsByKeyPaginated(ctx context.Context, key string, pageNum, limit int) ([]string, int64, error) {
	if r.client == nil {
		return nil, 0, fmt.Errorf("redis client is not initialized")
	}

	// Получаем общее количество элементов
	total, err := r.client.SCard(ctx, key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get suggestions count from Redis: %w", err)
	}

	// Если номер страницы некорректный или limit <= 0
	if pageNum < 0 || limit <= 0 {
		return []string{}, total, nil
	}

	// Вычисляем абсолютное смещение
	offset := pageNum * limit

	// Получаем все элементы
	allNames, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get suggestions from Redis: %w", err)
	}

	maxPage := int((total + int64(limit) - 1) / int64(limit)) // Округление вверх: (total + limit - 1) / limit

	// Если номер страницы превышает максимально возможный
	if pageNum >= maxPage {
		return []string{}, total, nil
	}

	// Применяем пагинацию для обычного случая
	end := offset + limit
	if end > len(allNames) {
		end = len(allNames)
	}

	return allNames[offset:end], total, nil
}

// GetProductSuggestionsByCategoryPaginated получает продукты для категории с пагинацией
// pageNum - номер страницы (начиная с 0)
// limit - количество элементов на странице
func (r *SuggestionsRepository) GetProductSuggestionsByCategoryPaginated(ctx context.Context, categoryID string, pageNum, limit int) ([]string, int64, error) {
	return r.GetSuggestionsByKeyPaginated(ctx, categoryKey(categoryID), pageNum, limit)
}
