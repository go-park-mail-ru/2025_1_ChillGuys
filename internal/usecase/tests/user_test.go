package tests

import (
	"context"
	"errors"
	"testing"

	minioMocks "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/user"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUsecase_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

	// Create a valid UUID for testing
	testUUID := uuid.New()
	testUUIDStr := testUUID.String()

	tests := []struct {
		name          string
		ctx           context.Context
		mockSetup     func()
		expectedUser  *dto.UserDTO
		expectedToken string
		expectedError error
	}{
		{
			name: "success",
			ctx: context.WithValue(
				context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
				domains.RoleKey{}, "buyer",
			),
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						ID:          testUUID,
						Email:       "test@example.com",
						Name:        "Test",
						Surname:     null.StringFrom("User"),
						ImageURL:    null.StringFrom("image.jpg"),
						PhoneNumber: null.StringFrom("1234567890"),
						Role:        models.RoleBuyer,
					}, nil)
			},
			expectedUser: &dto.UserDTO{
				ID:          testUUID,
				Email:       "test@example.com",
				Name:        "Test",
				Surname:     null.StringFrom("User"),
				ImageURL:    null.StringFrom("image.jpg"),
				PhoneNumber: null.StringFrom("1234567890"),
				Role:        "buyer",
			},
		},
		{
			name: "role changed - new token",
			ctx: context.WithValue(
				context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
				domains.RoleKey{}, "old-role",
			),
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						ID:   testUUID,
						Role: models.RoleSeller,
					}, nil)
				mockToken.EXPECT().
					CreateJWT(testUUIDStr, "seller").
					Return("new-token", nil)
			},
			expectedUser: &dto.UserDTO{
				ID:   testUUID,
				Role: "seller",
			},
			expectedToken: "new-token",
		},
		{
			name: "user not found",
			ctx: context.WithValue(
				context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
				domains.RoleKey{}, "buyer",
			),
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(nil, errs.ErrNotFound)
			},
			expectedError: errs.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			userDTO, token, err := uc.GetMe(tt.ctx)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedUser, userDTO)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestUserUsecase_UpdateUserProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

	// Create a valid UUID for testing
	testUUID := uuid.New()
	testUUIDStr := testUUID.String()

	tests := []struct {
		name          string
		ctx           context.Context
		update        dto.UpdateUserProfileRequestDTO
		mockSetup     func()
		expectedError error
	}{
		{
			name: "success - update all fields",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			update: dto.UpdateUserProfileRequestDTO{
				Name:        null.StringFrom("NewName"),
				Surname:     null.StringFrom("NewSurname"),
				PhoneNumber: null.StringFrom("1234567890"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						Name:        "OldName",
						Surname:     null.StringFrom("OldSurname"),
						PhoneNumber: null.StringFrom("111111"),
					}, nil)
				mockRepo.EXPECT().
					UpdateUserProfile(gomock.Any(), testUUID, models.UpdateUserDB{
						Name:        "NewName",
						Surname:     null.StringFrom("NewSurname"),
						PhoneNumber: null.StringFrom("1234567890"),
					}).
					Return(nil)
			},
		},
		{
			name: "partial update",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			update: dto.UpdateUserProfileRequestDTO{
				Name: null.StringFrom("NewName"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						Name:        "OldName",
						Surname:     null.StringFrom("OldSurname"),
						PhoneNumber: null.StringFrom("111111"),
					}, nil)
				mockRepo.EXPECT().
					UpdateUserProfile(gomock.Any(), testUUID, models.UpdateUserDB{
						Name:        "NewName",
						Surname:     null.StringFrom("OldSurname"),
						PhoneNumber: null.StringFrom("111111"),
					}).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := uc.UpdateUserProfile(tt.ctx, tt.update)
			if tt.expectedError != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestUserUsecase_BecomeSeller(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

	// Create a valid UUID for testing
	testUUID := uuid.New()
	testUUIDStr := testUUID.String()

	tests := []struct {
		name          string
		ctx           context.Context
		request       dto.UpdateRoleRequest
		mockSetup     func()
		expectedError error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			request: dto.UpdateRoleRequest{
				Title:       "Test Seller",
				Description: "Test Description",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateSellerAndUpdateRole(gomock.Any(), testUUID, "Test Seller", "Test Description").
					Return(nil)
			},
		},
		{
			name: "repository error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			request: dto.UpdateRoleRequest{
				Title:       "Test Seller",
				Description: "Test Description",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateSellerAndUpdateRole(gomock.Any(), testUUID, "Test Seller", "Test Description").
					Return(errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := uc.BecomeSeller(tt.ctx, tt.request)
			if tt.expectedError != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}