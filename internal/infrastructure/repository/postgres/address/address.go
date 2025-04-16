package address

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

const (
	queryCheckAddressExists = `
        SELECT id FROM bazaar.address
        WHERE address_string = $1 AND coordinate = $2
        LIMIT 1
    `
	queryUpsertAddress = `
        INSERT INTO bazaar.address (id, region, city, address_string, coordinate) 
        VALUES ($1, $2, $3, $4, $5)
    `
	queryUpsertUserAddress = `
        INSERT INTO bazaar.user_address (id, label, user_id, address_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, address_id) DO NOTHING
    `
	queryGetAddressesByUserID = `
		SELECT a.id, ua.label, a.region, a.city, a.address_string, a.coordinate
		FROM bazaar.address AS a
		JOIN bazaar.user_address AS ua ON a.id = ua.address_id
		WHERE ua.user_id = $1
	`
	queryGetAllPickupPoints = `
		SELECT a.id, a.region, a.city, a.address_string, a.coordinate
		FROM bazaar.pickup_point AS pp
		JOIN bazaar.address AS a ON pp.address_id = a.id
	`
)

//go:generate mockgen -source=address.go -destination=../mocks/address_repository_mock.go -package=mocks IAddressRepository
type IAddressRepository interface {
	CheckAddressExists(context.Context, models.AddressDB) (uuid.UUID, error)
	CreateAddress(context.Context, models.AddressDB) error
	CreateUserAddress(context.Context, models.UserAddress) error
	GetUserAddress(context.Context, uuid.UUID) (*[]dto.AddressDTO, error)
	GetAllPickupPoints(ctx context.Context) (*[]models.AddressDB, error)
}

type AddressRepository struct {
	db  *sql.DB
}

func NewAddressRepository(db *sql.DB) *AddressRepository {
	return &AddressRepository{
		db:  db,
	}
}

func (r *AddressRepository) CheckAddressExists(ctx context.Context, address models.AddressDB) (uuid.UUID, error) {
	const op = "AddressRepository.CheckAddressExists"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("address", address)
	
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, queryCheckAddressExists,
		address.AddressString, address.Coordinate,
	).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Debug("address not found")
			return uuid.Nil, nil
		}
		logger.WithError(err).Error("check address exists")
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *AddressRepository) CreateAddress(ctx context.Context, in models.AddressDB) error {
	const op = "AddressRepository.CreateAddress"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("address", in)

	_, err := r.db.QueryContext(ctx, queryUpsertAddress,
		in.ID.String(),
		in.Region,
		in.City,
		in.AddressString,
		in.Coordinate,
	)

	if err != nil {
		logger.WithError(err).Error("create address")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AddressRepository) CreateUserAddress(ctx context.Context, in models.UserAddress) error {
	const op = "AddressRepository.CreateUserAddress"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_address", in)

	if _, err := r.db.ExecContext(ctx, queryUpsertUserAddress,
		in.ID.String(),
		in.Label,
		in.UserID.String(),
		in.AddressID.String(),
	); err != nil {
		logger.WithError(err).Error("create user address")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AddressRepository) GetUserAddress(ctx context.Context, userID uuid.UUID) (*[]dto.AddressDTO, error) {
	const op = "AddressRepository.GetUserAddress"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	rows, err := r.db.QueryContext(ctx, queryGetAddressesByUserID, userID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("no addresses found for user")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("query user addresses")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var addresses []dto.AddressDTO

	for rows.Next() {
		var address dto.AddressDTO
		if err := rows.Scan(
			&address.ID,
			&address.Label,
			&address.Region,
			&address.City,
			&address.AddressString,
			&address.Coordinate,
		); err != nil {
			logger.WithError(err).Error("scan address row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, err
	}

	return &addresses, nil
}

func (r *AddressRepository) GetAllPickupPoints(ctx context.Context) (*[]models.AddressDB, error) {
	const op = "AddressRepository.GetAllPickupPoints"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	rows, err := r.db.QueryContext(ctx, queryGetAllPickupPoints)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("no pickup points found")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("query pickup points")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var points []models.AddressDB

	for rows.Next() {
		var addr models.AddressDB
		if err := rows.Scan(
			&addr.ID,
			&addr.Region,
			&addr.City,
			&addr.AddressString,
			&addr.Coordinate,
		); err != nil {
			logger.WithError(err).Error("scan pickup point row")
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		points = append(points, addr)
	}

	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("rows iteration error")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &points, nil
}
