package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckAddressExists(t *testing.T) {
	// Create test UUIDs once at the start
	existingAddressID := uuid.New()

	tests := []struct {
		name          string
		address       models.AddressDB
		mockBehavior  func(*mocks.MockIAddressRepository, models.AddressDB, uuid.UUID)
		expectedID    uuid.UUID
		expectedError error
	}{
		{
			name: "Success - address exists",
			address: models.AddressDB{
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Address"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, addr models.AddressDB, id uuid.UUID) {
				m.EXPECT().CheckAddressExists(gomock.Any(), addr).Return(id, nil)
			},
			expectedID:    existingAddressID, // Use the pre-generated UUID
			expectedError: nil,
		},
		{
			name: "Success - address not exists",
			address: models.AddressDB{
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Address"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, addr models.AddressDB, _ uuid.UUID) {
				m.EXPECT().CheckAddressExists(gomock.Any(), addr).Return(uuid.Nil, nil)
			},
			expectedID:    uuid.Nil,
			expectedError: nil,
		},
		{
			name: "Error - database error",
			address: models.AddressDB{
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Address"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, addr models.AddressDB, _ uuid.UUID) {
				m.EXPECT().CheckAddressExists(gomock.Any(), addr).Return(uuid.Nil, errors.New("database error"))
			},
			expectedID:    uuid.Nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockIAddressRepository(ctrl)
			test.mockBehavior(mockRepo, test.address, test.expectedID)

			id, err := mockRepo.CheckAddressExists(context.Background(), test.address)

			assert.Equal(t, test.expectedID, id)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateAddress(t *testing.T) {
	tests := []struct {
		name          string
		address       models.AddressDB
		mockBehavior  func(*mocks.MockIAddressRepository, models.AddressDB)
		expectedError error
	}{
		{
			name: "Success",
			address: models.AddressDB{
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Address"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, addr models.AddressDB) {
				m.EXPECT().CreateAddress(gomock.Any(), addr).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error",
			address: models.AddressDB{
				Region:        null.StringFrom("Region"),
				City:          null.StringFrom("City"),
				AddressString: null.StringFrom("Address"),
				Coordinate:    null.StringFrom("0,0"),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, addr models.AddressDB) {
				m.EXPECT().CreateAddress(gomock.Any(), addr).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockIAddressRepository(ctrl)
			test.mockBehavior(mockRepo, test.address)

			err := mockRepo.CreateAddress(context.Background(), test.address)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestCreateUserAddress(t *testing.T) {
	tests := []struct {
		name          string
		userAddress   models.UserAddress
		mockBehavior  func(*mocks.MockIAddressRepository, models.UserAddress)
		expectedError error
	}{
		{
			name: "Success",
			userAddress: models.UserAddress{
				ID:        uuid.New(),
				Label:     null.StringFrom("Home"), // Changed from sql.NullString to null.String
				UserID:    uuid.New(),
				AddressID: uuid.New(),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, ua models.UserAddress) {
				m.EXPECT().CreateUserAddress(gomock.Any(), ua).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Error",
			userAddress: models.UserAddress{
				ID:        uuid.New(),
				Label:     null.StringFrom("Home"),
				UserID:    uuid.New(),
				AddressID: uuid.New(),
			},
			mockBehavior: func(m *mocks.MockIAddressRepository, ua models.UserAddress) {
				m.EXPECT().CreateUserAddress(gomock.Any(), ua).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockIAddressRepository(ctrl)
			test.mockBehavior(mockRepo, test.userAddress)

			err := mockRepo.CreateUserAddress(context.Background(), test.userAddress)

			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetUserAddress(t *testing.T) {
	userID := uuid.New()
	testAddresses := []dto.AddressDTO{
		{
			ID:            uuid.New(),
			Label:         null.StringFrom("Home"),
			Region:        null.StringFrom("Region"),
			City:          null.StringFrom("City"),
			AddressString: null.StringFrom("Address"),
			Coordinate:    null.StringFrom("0,0"),
		},
		{
			ID:            uuid.New(),
			Label:         null.StringFrom("Work"),
			Region:        null.StringFrom("Region"),
			City:          null.StringFrom("City"),
			AddressString: null.StringFrom("Address2"),
			Coordinate:    null.StringFrom("1,1"),
		},
	}

	tests := []struct {
		name           string
		userID         uuid.UUID
		mockBehavior   func(*mocks.MockIAddressRepository, uuid.UUID)
		expectedResult *[]dto.AddressDTO
		expectedError  error
	}{
		{
			name:   "Success",
			userID: userID,
			mockBehavior: func(m *mocks.MockIAddressRepository, uid uuid.UUID) {
				m.EXPECT().GetUserAddress(gomock.Any(), uid).Return(&testAddresses, nil)
			},
			expectedResult: &testAddresses,
			expectedError:  nil,
		},
		{
			name:   "Error",
			userID: userID,
			mockBehavior: func(m *mocks.MockIAddressRepository, uid uuid.UUID) {
				m.EXPECT().GetUserAddress(gomock.Any(), uid).Return(nil, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockIAddressRepository(ctrl)
			test.mockBehavior(mockRepo, test.userID)

			result, err := mockRepo.GetUserAddress(context.Background(), test.userID)

			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func TestGetAllPickupPoints(t *testing.T) {
	testPoints := []models.AddressDB{
		{
			ID:            uuid.New(),
			Region:        null.StringFrom("Region1"),
			City:          null.StringFrom("City1"),
			AddressString: null.StringFrom("Address1"),
			Coordinate:    null.StringFrom("0,0"),
		},
		{
			ID:            uuid.New(),
			Region:        null.StringFrom("Region2"),
			City:          null.StringFrom("City2"),
			AddressString: null.StringFrom("Address2"),
			Coordinate:    null.StringFrom("1,1"),
		},
	}

	tests := []struct {
		name           string
		mockBehavior   func(*mocks.MockIAddressRepository)
		expectedResult *[]models.AddressDB
		expectedError  error
	}{
		{
			name: "Success",
			mockBehavior: func(m *mocks.MockIAddressRepository) {
				m.EXPECT().GetAllPickupPoints(gomock.Any()).Return(&testPoints, nil)
			},
			expectedResult: &testPoints,
			expectedError:  nil,
		},
		{
			name: "Error",
			mockBehavior: func(m *mocks.MockIAddressRepository) {
				m.EXPECT().GetAllPickupPoints(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mocks.NewMockIAddressRepository(ctrl)
			test.mockBehavior(mockRepo)

			result, err := mockRepo.GetAllPickupPoints(context.Background())

			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
