package address

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IAddressUsecase interface {
	CreateAddress(ctx context.Context, in models.Address) error
	GetAddresses(context.Context, uuid.UUID) ([]models.Address, error)
	GetPickupPoints(ctx context.Context) ([]models.AddressDB, error)
}

type AddressUsecase struct {
	repo address.IAddressRepository
	log  *logrus.Logger
}

func NewAddressUsecase(
	repo address.IAddressRepository,
	log *logrus.Logger,
) *AddressUsecase {
	return &AddressUsecase{
		repo: repo,
		log:  log,
	}
}

func (u *AddressUsecase) CreateAddress(ctx context.Context, in models.Address) error {
	userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
	if !isExist {
		return errs.ErrNotFound
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errs.ErrInvalidID
	}

	addressID := uuid.New()
	addr := models.AddressDB{
		ID:        addressID,
		City:      in.City,
		Street:    in.Street,
		House:     in.House,
		Apartment: in.Apartment,
		ZipCode:   in.ZipCode,
	}

	addrID, err := u.repo.CheckAddressExists(ctx, addr)
	if err != nil {
		return err
	}

	// Если адреса не существует, мы его создаём
	if addrID == uuid.Nil {
		if err = u.repo.CreateAddress(ctx, addr); err != nil {
			return err
		}
	} else {
		addressID = addrID
	}

	// Запоминает адрес в таблице адресов пользователя
	userAddr := models.UserAddress{
		ID:        uuid.New(),
		Label:     in.Label,
		UserID:    userID,
		AddressID: addressID,
	}

	return u.repo.CreateUserAddress(ctx, userAddr)
}

func (u *AddressUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]models.Address, error) {
	addresses, err := u.repo.GetUserAddress(ctx, userID)
	if err != nil {
		return nil, err
	}

	return *addresses, nil
}

func (u *AddressUsecase) GetPickupPoints(ctx context.Context) ([]models.AddressDB, error) {
	points, err := u.repo.GetAllPickupPoints(ctx)
	if err != nil {
		return nil, err
	}

	var res []models.AddressDB
	for _, point := range *points {
		res = append(res, models.AddressDB{
			ID:        point.ID,
			City:      point.City,
			Street:    point.Street,
			House:     point.House,
			Apartment: point.Apartment,
			ZipCode:   point.ZipCode,
		})
	}

	return res, nil
}
