package address

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	queryCheckAddressExists = `
        SELECT id FROM bazaar.address
        WHERE region = $1 AND city = $2 AND address_string = $3 AND coordinate = $4
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

type IAddressRepository interface {
	CheckAddressExists(context.Context, models.AddressDB) (uuid.UUID, error)
	CreateAddress(context.Context, models.AddressDB) error
	CreateUserAddress(context.Context, models.UserAddress) error
	GetUserAddress(context.Context, uuid.UUID) (*[]models.Address, error)
	GetAllPickupPoints(ctx context.Context) (*[]models.AddressDB, error)
}

type AddressRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewAddressRepository(db *sql.DB, log *logrus.Logger) *AddressRepository {
	return &AddressRepository{
		db:  db,
		log: log,
	}
}

func (r *AddressRepository) CheckAddressExists(ctx context.Context, address models.AddressDB) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, queryCheckAddressExists,
		address.Region, address.City, address.AddressString, address.Coordinate,
	).Scan(&id)

	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	} else if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *AddressRepository) CreateAddress(ctx context.Context, in models.AddressDB) error {
	_, err := r.db.QueryContext(ctx, queryUpsertAddress,
		in.ID.String(),
		in.Region,
		in.City,
		in.AddressString,
		in.Coordinate,
	)

	return err
}

func (r *AddressRepository) CreateUserAddress(ctx context.Context, in models.UserAddress) error {
	if _, err := r.db.ExecContext(ctx, queryUpsertUserAddress,
		in.ID.String(),
		in.Label,
		in.UserID.String(),
		in.AddressID.String(),
	); err != nil {
		return err
	}

	return nil
}

func (r *AddressRepository) GetUserAddress(ctx context.Context, userID uuid.UUID) (*[]models.Address, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAddressesByUserID, userID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.Address

	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.Label,
			&address.Region,
			&address.City,
			&address.AddressString,
			&address.Coordinate,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &addresses, nil
}

func (r *AddressRepository) GetAllPickupPoints(ctx context.Context) (*[]models.AddressDB, error) {
	rows, err := r.db.QueryContext(ctx, queryGetAllPickupPoints)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		points = append(points, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &points, nil
}
