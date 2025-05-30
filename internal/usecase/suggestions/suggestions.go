package suggestions

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"strings"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

//go:generate mockgen -source=suggestions.go -destination=../../infrastructure/repository/postgres/mocks/suggestions_repository_mock.go -package=mocks ISuggestionsRepository
type ISuggestionsRepository interface {
	GetAllCategoriesName(ctx context.Context) ([]*models.CategorySuggestion, error)
	GetAllProductsName(ctx context.Context) ([]*models.ProductSuggestion, error)
	GetProductsNameByCategory(ctx context.Context, categoryID string) ([]*models.ProductSuggestion, error)
}

type ISuggestionsRedisRepository interface {
	AddProductSuggestionsByCategory(ctx context.Context, categoryID string, names []string) error
	AddSuggestionsByKey(ctx context.Context, key string, names []string) error
	GetProductSuggestionsByCategory(ctx context.Context, categoryID string) ([]string, error)
	GetProductSuggestionsByCategoryPaginated(ctx context.Context, categoryID string, pageNum int, limit int) ([]string, int64, error)
	GetSuggestionsByKey(ctx context.Context, key string) ([]string, error)
}

type SuggestionsUsecase struct {
	repo      ISuggestionsRepository
	redisRepo ISuggestionsRedisRepository
}

func NewSuggestionsUsecase(repo ISuggestionsRepository, redisRepo ISuggestionsRedisRepository) *SuggestionsUsecase {
	return &SuggestionsUsecase{
		repo:      repo,
		redisRepo: redisRepo,
	}
}

func (u *SuggestionsUsecase) GetCategorySuggestions(ctx context.Context, subString string) (dto.CategoryNameResponse, error) {
	const op = "SuggestionsUsecase.GetCategorySuggestions"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	names, err := u.redisRepo.GetSuggestionsByKey(ctx, redis.CategoryNamesKey)
	if err != nil {
		logger.WithError(err).Error("failed to get category suggestions from Redis")
	}

	if len(names) == 0 {
		categories, err := u.repo.GetAllCategoriesName(ctx)
		if err != nil {
			logger.WithError(err).Error("get categories from repository")
			return dto.CategoryNameResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		names = make([]string, 0, len(categories))
		for _, c := range categories {
			names = append(names, c.Name)
		}
		_ = u.redisRepo.AddSuggestionsByKey(ctx, redis.CategoryNamesKey, names)
	}

	filtered := filterSuggestions(names, subString)

	var suggestions []models.CategorySuggestion
	for _, name := range filtered {
		suggestions = append(suggestions, models.CategorySuggestion{Name: name})
	}

	return dto.CategoryNameResponse{CategoriesNames: suggestions}, nil
}

func (u *SuggestionsUsecase) GetProductSuggestions(ctx context.Context, categoryID null.String, subString string) (dto.ProductNameResponse, error) {
	const op = "SuggestionsUsecase.GetProductSuggestions"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var names []string
	var err error

	// Пробуем достать из Redis
	if categoryID.Valid {
		names, err = u.redisRepo.GetProductSuggestionsByCategory(ctx, categoryID.String)
		if err != nil {
			logger.WithError(err).Error("failed to get category-specific suggestions from Redis")
		}
	} else {
		names, err = u.redisRepo.GetSuggestionsByKey(ctx, redis.ProductNamesKey)
		if err != nil {
			logger.WithError(err).Error("failed to get general product suggestions from Redis")
		}
	}

	// Если в Redis ничего нет — грузим из БД
	if len(names) == 0 {
		var products []*models.ProductSuggestion

		if categoryID.Valid {
			products, err = u.repo.GetProductsNameByCategory(ctx, categoryID.String)
			if err != nil {
				logger.WithError(err).Error("get products by category from repository")
				return dto.ProductNameResponse{}, fmt.Errorf("%s: %w", op, err)
			}
		} else {
			products, err = u.repo.GetAllProductsName(ctx)
			if err != nil {
				logger.WithError(err).Error("get all products from repository")
				return dto.ProductNameResponse{}, fmt.Errorf("%s: %w", op, err)
			}
		}

		names = make([]string, 0, len(products))
		for _, p := range products {
			names = append(names, p.Name)
		}

		// Кэшируем в Redis
		if categoryID.Valid {
			_ = u.redisRepo.AddProductSuggestionsByCategory(ctx, categoryID.String, names)
		} else {
			_ = u.redisRepo.AddSuggestionsByKey(ctx, redis.ProductNamesKey, names)
		}
	}

	// Фильтрация по подстроке
	filtered := filterSuggestions(names, subString)

	// Возврат
	suggestions := make([]models.ProductSuggestion, 0, len(filtered))
	for _, name := range filtered {
		suggestions = append(suggestions, models.ProductSuggestion{Name: name})
	}

	return dto.ProductNameResponse{ProductNames: suggestions}, nil
}

func filterSuggestions(names []string, sub string) []string {
	lowerSub := strings.ToLower(sub)
	seen := make(map[string]struct{})
	var result []string

	for _, fullName := range names {
		lowerFull := strings.ToLower(fullName)

		if strings.HasPrefix(lowerFull, lowerSub) {
			if _, ok := seen[fullName]; !ok {
				result = append(result, fullName)
				seen[fullName] = struct{}{}
			}
			continue
		}

		for _, word := range strings.Fields(lowerFull) {
			if strings.HasPrefix(word, lowerSub) {
				if _, ok := seen[fullName]; !ok {
					result = append(result, fullName)
					seen[fullName] = struct{}{}
				}
				break
			}
		}
	}

	return result
}
