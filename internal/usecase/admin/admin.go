package admin

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	productRepo "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/product"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

//go:generate mockgen -source=admin.go -destination=../../infrastructure/repository/postgres/mocks/admin_repository_mock.go -package=mocks IAdminRepository
type IAdminRepository interface {
	GetPendingProducts(ctx context.Context, offset int) ([]*models.Product, error)
	UpdateProductStatus(ctx context.Context, productID uuid.UUID, status models.ProductStatus) error
	GetPendingUsers(ctx context.Context, offset int) ([]*models.User, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error
}

type AdminUsecase struct {
	repo        IAdminRepository
	repoProduct productRepo.IProductRepository
	redisRepo   *redis.SuggestionsRepository
}

func NewAdminUsecase(r IAdminRepository, redisRepo *redis.SuggestionsRepository, repoProduct productRepo.IProductRepository) *AdminUsecase {
	return &AdminUsecase{
		repo:        r,
		repoProduct: repoProduct,
		redisRepo:   redisRepo,
	}
}

func (u *AdminUsecase) GetPendingProducts(ctx context.Context, offset int) (dto.ProductsResponse, error) {
	const op = "AdminUsecase.GetPendingProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	products, err := u.repo.GetPendingProducts(ctx, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get pending products")
		return dto.ProductsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ConvertToProductsResponse(products), nil
}

func (u *AdminUsecase) UpdateProductStatus(ctx context.Context, req dto.UpdateProductStatusRequest) error {
	const op = "AdminUsecase.UpdateProductStatus"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_id", req.ProductID)

	var status models.ProductStatus
	switch req.Update {
	case 0:
		status = models.ProductRejected
	case 1:
		status = models.ProductApproved
	default:
		logger.Error("invalid update status value")
		return errs.ErrParseRequestData
	}

	err := u.repo.UpdateProductStatus(ctx, req.ProductID, status)
	if err != nil {
		logger.WithError(err).Error("failed to update product status")
		return fmt.Errorf("%s: %w", op, err)
	}

	product, err := u.repoProduct.GetProductByID(ctx, req.ProductID)
	if err != nil {
		logger.WithError(err).Error("failed to get product details")
		return fmt.Errorf("%s: %w", op, err)
	}

	productNames := []string{product.Name}

	err = u.redisRepo.AddSuggestionsByKey(ctx, redis.ProductNamesKey, productNames)
	if err != nil {
		logger.WithError(err).Error("failed to add product to Redis suggestions")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *AdminUsecase) GetPendingUsers(ctx context.Context, offset int) (dto.UsersResponse, error) {
	const op = "AdminUsecase.GetPendingUsers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	users, err := u.repo.GetPendingUsers(ctx, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get pending users")
		return dto.UsersResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.ConvertToUsersResponse(users), nil
}

func (u *AdminUsecase) UpdateUserRole(ctx context.Context, req dto.UpdateUserRoleRequest) error {
	const op = "AdminUsecase.UpdateUserRole"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", req.UserID)

	var role models.UserRole
	switch req.Update {
	case 0:
		role = models.RoleBuyer
	case 1:
		role = models.RoleSeller
	default:
		logger.Error("invalid update role value")
		return errs.ErrParseRequestData
	}

	err := u.repo.UpdateUserRole(ctx, req.UserID, role)
	if err != nil {
		logger.WithError(err).Error("failed to update user role")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
