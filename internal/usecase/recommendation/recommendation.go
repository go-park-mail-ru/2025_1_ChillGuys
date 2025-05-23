package recommendation

import (
	"context"
	"fmt"
	recommendationRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/recommendation"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/product"
	"github.com/google/uuid"
	"sync"
)

type IRecommendationUsecase interface {
	GetRecommendations(ctx context.Context, productID uuid.UUID) ([]*models.Product, error)
}

type RecommendationUsecase struct {
	pu product.IProductUsecase
	rr recommendationRepo.IRecommendationRepository
}

func NewRecommendationUsecase(
	productUscase product.IProductUsecase,
	recommendationRepo recommendationRepo.IRecommendationRepository,
) *RecommendationUsecase {
	return &RecommendationUsecase{
		rr: recommendationRepo,
		pu: productUscase,
	}
}

func (u *RecommendationUsecase) GetRecommendations(ctx context.Context, productID uuid.UUID) ([]*models.Product, error) {
	const op = "RecommendationUsecase.GetRecommendations"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	subcatIDs, err := u.rr.GetCategoryIDsByProductID(ctx, productID)
	if err != nil {
		logger.WithError(err).Error("get subcategory ids from repository")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var productIDsByCategory []uuid.UUID
	var productIDRec []uuid.UUID
	savedProductIDs := make(map[uuid.UUID]bool)

	for _, subcatID := range subcatIDs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ctx.Err() != nil {
				return
			}

			productIDs, err := u.rr.GetProductIDsBySubcategoryID(ctx, subcatID, 10)
			if err != nil {
				logger.WithError(err).WithField("product_ids", productIDs).Warn("failed to get product IDs by Subcategory IDs")
				return
			}

			mu.Lock()
			for _, id := range productIDs {
				if _, ok := savedProductIDs[id]; !ok {
					productIDsByCategory = append(productIDsByCategory, id)
					savedProductIDs[id] = true
				}
			}
			mu.Unlock()
		}()
	}

	wg.Wait()

	if len(productIDsByCategory) == 0 {
		return nil, nil
	}

	productIDsByCategory = productIDsByCategory[:min(10, len(productIDsByCategory))]

	for _, id := range productIDsByCategory {
		if id != productID {
			productIDRec = append(productIDRec, id)
		}
	}

	productsRecommendation, err := u.pu.GetProductsByIDs(ctx, productIDRec)
	if err != nil {
		logger.WithError(err).Error("get productsRecommendation by IDs")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return productsRecommendation, nil
}
