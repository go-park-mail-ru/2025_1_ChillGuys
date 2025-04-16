package address

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=address.go -destination=../mocks/address_usecase_mock.go -package=mocks IAddressUsecase
type IAddressUsecase interface {
	CreateAddress(context.Context, uuid.UUID, dto.AddressDTO) error
	GetAddresses(context.Context, uuid.UUID) ([]dto.GetAddressResDTO, error)
	GetPickupPoints(ctx context.Context) ([]dto.GetPointAddressResDTO, error)
}

type AddressUsecase struct {
	Repo address.IAddressRepository
	Log  *logrus.Logger
}

func NewAddressUsecase(
	repo address.IAddressRepository,
	log *logrus.Logger,
) *AddressUsecase {
	return &AddressUsecase{
		Repo: repo,
		Log:  log,
	}
}

func (u *AddressUsecase) CreateAddress(ctx context.Context, userID uuid.UUID, in dto.AddressDTO) error {
	addressID := uuid.New()
	addr := models.AddressDB{
		ID:            addressID,
		Region:        in.Region,
		City:          in.City,
		AddressString: in.AddressString,
		Coordinate:    in.Coordinate,
	}

	addrID, err := u.Repo.CheckAddressExists(ctx, addr)
	if err != nil {
		return err
	}

	if addrID == uuid.Nil {
		if err = u.Repo.CreateAddress(ctx, addr); err != nil {
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

	return u.Repo.CreateUserAddress(ctx, userAddr)
}

func (u *AddressUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]dto.GetAddressResDTO, error) {
	addresses, err := u.Repo.GetUserAddress(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.GetAddressResDTO, 0, len(*addresses))
	for _, addr := range *addresses {
		res = append(res, dto.GetAddressResDTO{
			ID:            addr.ID,
			Label:         addr.Label,
			AddressString: addr.AddressString,
			Coordinate:    addr.Coordinate,
		})
	}

	return res, nil
}

func (u *AddressUsecase) GetPickupPoints(ctx context.Context) ([]dto.GetPointAddressResDTO, error) {
	points, err := u.Repo.GetAllPickupPoints(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.GetPointAddressResDTO, 0, len(*points))
	for _, point := range *points {
		res = append(res, dto.GetPointAddressResDTO{
			ID:            point.ID,
			AddressString: point.AddressString,
			Coordinate:    point.Coordinate,
		})
	}

	return res, nil
}
