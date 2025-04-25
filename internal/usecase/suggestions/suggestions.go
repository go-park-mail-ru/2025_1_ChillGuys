package suggestions

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

type ISuggestionsRepository interface {
	GetAllCategoriesName(ctx context.Context) ([]*models.CategorySuggestion, error)
	GetAllProductsName(ctx context.Context) ([]*models.ProductSuggestion, error)
}

type SuggestionsUsecase struct {
	repo      ISuggestionsRepository
	redisRepo *redis.SuggestionsRepository
}

func NewSuggestionsUsecase(repo ISuggestionsRepository, redisRepo *redis.SuggestionsRepository) *SuggestionsUsecase {
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

func (u *SuggestionsUsecase) GetProductSuggestions(ctx context.Context, subString string) (dto.ProductNameResponse, error) {
	const op = "SuggestionsUsecase.GetProductSuggestions"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	names, err := u.redisRepo.GetSuggestionsByKey(ctx, redis.ProductNamesKey)
	if err != nil {
		logger.WithError(err).Error("failed to get product suggestions from Redis")
	}

	if len(names) == 0 {
		products, err := u.repo.GetAllProductsName(ctx)
		if err != nil {
			logger.WithError(err).Error("get products from repository")
			return dto.ProductNameResponse{}, fmt.Errorf("%s: %w", op, err)
		}

		names = make([]string, 0, len(products))
		for _, p := range products {
			names = append(names, p.Name)
		}

		_ = u.redisRepo.AddSuggestionsByKey(ctx, redis.ProductNamesKey, names)
	}

	filtered := filterSuggestions(names, subString)

	var suggestions []models.ProductSuggestion
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
