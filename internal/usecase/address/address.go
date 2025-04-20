package address

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/address"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/google/uuid"
)

//go:generate mockgen -source=address.go -destination=../mocks/address_usecase_mock.go -package=mocks IAddressUsecase
type IAddressUsecase interface {
	CreateAddress(context.Context, uuid.UUID, dto.AddressDTO) error
	GetAddresses(context.Context, uuid.UUID) ([]dto.GetAddressResDTO, error)
	GetPickupPoints(ctx context.Context) ([]dto.GetPointAddressResDTO, error)
}

type AddressUsecase struct {
	Repo address.IAddressRepository
}

func NewAddressUsecase(
	repo address.IAddressRepository,
) *AddressUsecase {
	return &AddressUsecase{
		Repo: repo,
	}
}

func (u *AddressUsecase) CreateAddress(ctx context.Context, userID uuid.UUID, in dto.AddressDTO) error {
	const op = "AddressUsecase.CreateAddress"
	logger := logctx.GetLogger(ctx).WithField("op", op).
		WithField("user_id", userID).
		WithField("address", in)

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
			logger.WithError(err).Error("create address")
			return fmt.Errorf("%s: %w", op, err)
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

	if err := u.Repo.CreateUserAddress(ctx, userAddr); err != nil {
		logger.WithError(err).Error("create user address")
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *AddressUsecase) GetAddresses(ctx context.Context, userID uuid.UUID) ([]dto.GetAddressResDTO, error) {
	const op = "AddressUsecase.GetAddresses"
	logger := logctx.GetLogger(ctx).WithField("op", op).WithField("user_id", userID)

	addresses, err := u.Repo.GetUserAddress(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("no addresses found for user")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("get user addresses")
		return nil, fmt.Errorf("%s: %w", op, err)
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
	const op = "AddressUsecase.GetPickupPoints"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	points, err := u.Repo.GetAllPickupPoints(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			logger.Warn("no pickup points found")
			return nil, fmt.Errorf("%s: %w", op, errs.NewNotFoundError(op))
		}
		logger.WithError(err).Error("get pickup points")
		return nil, fmt.Errorf("%s: %w", op, err)
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
