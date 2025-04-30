package admin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryGetPendingProducts = `
		SELECT 
			p.id, 
			p.seller_id, 
			p.name, 
			p.preview_image_url, 
			p.description, 
			p.status, 
			p.price, 
			p.quantity, 
			p.updated_at, 
			p.rating, 
			p.reviews_count,
			d.discounted_price
		FROM 
			bazaar.product p
		LEFT JOIN 
			bazaar.discount d ON p.id = d.product_id
		WHERE 
			p.status = 'pending'
		LIMIT 20 OFFSET $1`

	queryUpdateStatusProduct = `
		UPDATE bazaar.product
		SET 
			status = $1,
			updated_at = now()
		WHERE 
			id = $2`

	queryGetPendingUsers = `
	SELECT 
		u.id,
		u.email,
		u.name,
		u.surname,
		u.image_url,
		u.phone_number,
		u.role,
		s.id,
		s.title,
		s.description,
		s.user_id
	FROM 
		bazaar."user" u
	LEFT JOIN 
		bazaar.seller s ON u.id = s.user_id
	WHERE 
		u.role = 'pending'
	LIMIT 20 OFFSET $1`

	queryUpdateRoleUser = `
		UPDATE bazaar."user"
		SET 
			role = $1
		WHERE 
			id = $2`
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetPendingProducts возвращает список товаров со статусом "pending" с пагинацией
func (r *AdminRepository) GetPendingProducts(ctx context.Context, offset int) ([]*models.Product, error) {
	const op = "AdminRepository.GetPendingProducts"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	products := make([]*models.Product, 0)

	rows, err := r.db.QueryContext(ctx, queryGetPendingProducts, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query pending products")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var product models.Product
		var priceDiscount sql.NullFloat64
		err := rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.PreviewImageURL,
			&product.Description,
			&product.Status,
			&product.Price,
			&product.Quantity,
			&product.UpdatedAt,
			&product.Rating,
			&product.ReviewsCount,
			&priceDiscount,
		)
		if err != nil {
			logger.WithError(err).Error("failed to scan product row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		product.PriceDiscount = priceDiscount.Float64
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}

// UpdateProductStatus обновляет статус товара и возвращает обновленный товар
func (r *AdminRepository) UpdateProductStatus(ctx context.Context, productID uuid.UUID, status models.ProductStatus) error {
	const op = "AdminRepository.UpdateProductStatus"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("product_id", productID)

	_, err := r.db.ExecContext(ctx, queryUpdateStatusProduct, status.String(), productID)
	if err != nil {
		logger.WithError(err).Error("update product status")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetPendingUsers возвращает список пользователей с ролью "pending" с пагинацией
func (r *AdminRepository) GetPendingUsers(ctx context.Context, offset int) ([]*models.User, error)  {
	const op = "AdminRepository.GetPendingUsers"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	users := make([]*models.User, 0)

	rows, err := r.db.QueryContext(ctx, queryGetPendingUsers, offset)
	if err != nil {
		logger.WithError(err).Error("failed to query pending users")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
        var seller models.Seller
        var sellerID uuid.NullUUID

		err := rows.Scan(
            &user.ID,
            &user.Email,
            &user.Name,
            &user.Surname,
            &user.ImageURL,
            &user.PhoneNumber,
            &user.Role,
            &sellerID,
            &seller.Title,
            &seller.Description,
            &seller.UserID,
        )
		if err != nil {
			logger.WithError(err).Error("failed to scan user row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if sellerID.Valid {
            seller.ID = sellerID.UUID
            user.Seller = &seller
        }
		
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

// UpdateUserRole обновляет роль пользователя и возвращает обновленного пользователя
func (r *AdminRepository) UpdateUserRole(ctx context.Context, userID uuid.UUID, role models.UserRole) error {
	const op = "AdminRepository.UpdateUserRole"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	_, err := r.db.ExecContext(ctx, queryUpdateRoleUser, role, userID)
	if err != nil {
		logger.WithError(err).Error("update user role")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}