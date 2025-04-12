package address

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IAddressUsecase interface {
	CreateAddress(context.Context, uuid.UUID, models.Address) error
	GetAddresses(context.Context, uuid.UUID) ([]models.GetAddressRes, error)
	GetPickupPoints(ctx context.Context) ([]models.GetPointAddressRes, error)
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

func (u *AddressUsecase) CreateAddress(ctx context.Context, userID uuid.UUID, in models.Address) error {
	addressID := uuid.New()
	addr := models.AddressDB{
		ID:            addressID,
		Region:        in.Region,
		City:          in.City,
		AddressString: in.AddressString,
		Coordinate:    in.Coordinate,
	}

	addrID, err := u.repo.CheckAddressExists(ctx, addr)
	if err != nil {
		return err
	}

	if addrID == uuid.Nil {
		if err = u.repo.CreateAddress(ctx, addr); err != nil {
			return err
		}
	} else {
		addressID = addrID
	}

	userAddr := models.UserAddress{
		ID:        uuid.New(),
		Label:     in.Label,
		UserID:    userID,
		AddressID: addressID,
	}

	return u.repo.CreateUserAddress(ctx, userAddr)
}

func (u *AddressUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]models.GetAddressRes, error) {
	addresses, err := u.repo.GetUserAddress(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]models.GetAddressRes, 0, len(*addresses))
	for _, addr := range *addresses {
		res = append(res, models.GetAddressRes{
			ID:            addr.ID,
			Label:         addr.Label,
			AddressString: addr.AddressString,
			Coordinate:    addr.Coordinate,
		})
	}

	return res, nil
}

func (u *AddressUsecase) GetPickupPoints(ctx context.Context) ([]models.GetPointAddressRes, error) {
	points, err := u.repo.GetAllPickupPoints(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]models.GetPointAddressRes, 0, len(*points))
	for _, point := range *points {
		res = append(res, models.GetPointAddressRes{
			ID:            point.ID,
			AddressString: point.AddressString,
			Coordinate:    point.Coordinate,
		})
	}

	return res, nil
}
