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
        WHERE city = $1 AND street = $2 AND house = $3 AND apartment = $4 AND zip_code = $5
        LIMIT 1
    `
	queryUpsertAddress = `
        INSERT INTO bazaar.address (id, city, street, house, apartment, zip_code) 
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	queryUpsertUserAddress = `
        INSERT INTO bazaar.user_address (id, label, user_id, address_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id, address_id) DO NOTHING
    `
	queryGetAddressesByUserID = `
		SELECT a.id, ua.label, a.city, a.street, a.house, a.apartment, a.zip_code
		FROM bazaar.address AS a
		JOIN bazaar.user_address AS ua ON a.id = ua.address_id
		WHERE ua.user_id = $1
	`
	queryGetAllPickupPoints = `
		SELECT a.id, a.city, a.street, a.house, a.apartment, a.zip_code
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
		address.City, address.Street, address.House, address.Apartment, address.ZipCode,
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
		in.City,
		in.Street,
		in.House,
		in.Apartment,
		in.ZipCode,
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
		err := rows.Scan(
			&address.ID,
			&address.Label,
			&address.City,
			&address.Street,
			&address.House,
			&address.Apartment,
			&address.ZipCode,
		)
		if err != nil {
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
		err := rows.Scan(
			&addr.ID,
			&addr.City,
			&addr.Street,
			&addr.House,
			&addr.Apartment,
			&addr.ZipCode,
		)
		if err != nil {
			return nil, err
		}
		points = append(points, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &points, nil
}
