package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	addressus "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/address"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		userID        uuid.UUID
		input         dto.AddressDTO
		mockSetup     func(repo *mocks.MockIAddressRepository)
		expectedError error
	}{
		{
			name:   "Success - new address",
			userID: uuid.New(),
			input: dto.AddressDTO{
				Label:         null.StringFrom("Home"),
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Street 1"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				repo.EXPECT().
					CheckAddressExists(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, nil)

				repo.EXPECT().
					CreateAddress(gomock.Any(), gomock.Any()).
					Return(nil)

				repo.EXPECT().
					CreateUserAddress(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Success - existing address",
			userID: uuid.New(),
			input: dto.AddressDTO{
				Label:         null.StringFrom("Work"),
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Street 2"),
				Coordinate:    null.StringFrom("1,1"),
			},
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				existingID := uuid.New()
				repo.EXPECT().
					CheckAddressExists(gomock.Any(), gomock.Any()).
					Return(existingID, nil)

				repo.EXPECT().
					CreateUserAddress(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "Error - CheckAddressExists fails",
			userID: uuid.New(),
			input: dto.AddressDTO{
				Label:         null.StringFrom("Home"),
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Street 1"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				repo.EXPECT().
					CheckAddressExists(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockIAddressRepository(ctrl)
			tt.mockSetup(repo)

			uc := addressus.NewAddressUsecase(repo, logrus.New())
			err := uc.CreateAddress(context.Background(), tt.userID, tt.input)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAddresses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         uuid.UUID
		mockSetup      func(repo *mocks.MockIAddressRepository)
		expectedResult []dto.GetAddressResDTO
		expectedError  error
	}{
		{
			name:   "Success - get addresses",
			userID: uuid.New(),
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				mockResult := []dto.AddressDTO{
					{
						ID:            uuid.New(),
						Label:         null.StringFrom("Home"),
						Region:        null.StringFrom("Region1"),
						City:          null.StringFrom("City1"),
						AddressString: null.StringFrom("Street 1"),
						Coordinate:    null.StringFrom("0,0"),
					},
				}
				repo.EXPECT().
					GetUserAddress(gomock.Any(), gomock.Any()).
					Return(&mockResult, nil)
			},
			expectedResult: []dto.GetAddressResDTO{
				{
					Label:         null.StringFrom("Home"),
					AddressString: null.StringFrom("Street 1"),
					Coordinate:    null.StringFrom("0,0"),
				},
			},
			expectedError: nil,
		},
		{
			name:   "Error - repository fails",
			userID: uuid.New(),
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				repo.EXPECT().
					GetUserAddress(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockIAddressRepository(ctrl)
			tt.mockSetup(repo)

			uc := addressus.NewAddressUsecase(repo, logrus.New())
			res, err := uc.GetAddresses(context.Background(), tt.userID)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(res))
			}
		})
	}
}

func TestGetPickupPoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		mockSetup      func(repo *mocks.MockIAddressRepository)
		expectedResult []dto.GetPointAddressResDTO
		expectedError  error
	}{
		{
			name: "Success - get pickup points",
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				mockResult := []models.AddressDB{
					{
						ID:            uuid.New(),
						AddressString: null.StringFrom("Point 1"),
						Coordinate:    null.StringFrom("0,0"),
					},
				}
				repo.EXPECT().
					GetAllPickupPoints(gomock.Any()).
					Return(&mockResult, nil)
			},
			expectedResult: []dto.GetPointAddressResDTO{
				{
					AddressString: null.StringFrom("Point 1"),
					Coordinate:    null.StringFrom("0,0"),
				},
			},
			expectedError: nil,
		},
		{
			name: "Error - repository fails",
			mockSetup: func(repo *mocks.MockIAddressRepository) {
				repo.EXPECT().
					GetAllPickupPoints(gomock.Any()).
					Return(nil, errors.New("db error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockIAddressRepository(ctrl)
			tt.mockSetup(repo)

			uc := addressus.NewAddressUsecase(repo, logrus.New())
			res, err := uc.GetPickupPoints(context.Background())

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(res))
			}
		})
	}
}
