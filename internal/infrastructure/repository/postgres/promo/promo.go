package promo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
)

const (
	queryCreatePromoCode = `
		INSERT INTO bazaar.promo_code (id, code, percent, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
	`

	queryGetAllPromoCode = `
		SELECT id, code, percent, start_date, end_date
		FROM bazaar.promo_code
		ORDER BY start_date DESC
		LIMIT 50 OFFSET $1
	`

	queryCheckPromo = `
        SELECT id, code, percent, start_date, end_date 
        FROM bazaar.promo_code 
        WHERE code = $1
    `
)

type PromoRepository struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) *PromoRepository {
	return &PromoRepository{db:db}
}

func (r *PromoRepository) Create(ctx context.Context, promo models.PromoCode) error {
	const op = "PromoRepository.CreatePromoCode"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	_, err := r.db.ExecContext(ctx, queryCreatePromoCode,
		promo.ID,
		promo.Code,
		promo.Percent,
		promo.StartDate,
		promo.EndDate,
	)
	if err != nil {
		logger.WithError(err).Error("create promo code")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PromoRepository) GetAll(ctx context.Context, offset int) ([]*models.PromoCode, error) {
	const op = "PromoRepository.GetAllPromoCodes"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	promoList := []*models.PromoCode{}

	rows, err := r.db.QueryContext(ctx, queryGetAllPromoCode, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("no promo codes found")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("query promo codes")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		promo := &models.PromoCode{}
		err = rows.Scan(
			&promo.ID,
			&promo.Code,
			&promo.Percent,
			&promo.StartDate,
			&promo.EndDate,
		)
		if err != nil {
			logger.WithError(err).Error("scan promo code")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		promoList = append(promoList, promo)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return promoList, nil
}

func (r *PromoRepository) CheckPromoCode(ctx context.Context, code string) (*models.PromoCode, error) {
    const op = "PromoRepository.GetAllPromoCodes"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	
	var promo models.PromoCode
    err := r.db.QueryRowContext(ctx, queryCheckPromo, code).Scan(
        &promo.ID,
        &promo.Code,
        &promo.Percent,
        &promo.StartDate,
        &promo.EndDate,
    )
    
    if err != nil {
		logger.WithError(err).Error("scan promo code")
        return nil, err
    }
    
    return &promo, nil
}