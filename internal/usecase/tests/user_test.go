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
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_GetMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

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
		{
			name: "invalid user ID",
			ctx: context.WithValue(
				context.WithValue(context.Background(), domains.UserIDKey{}, "invalid"),
				domains.RoleKey{}, "buyer",
			),
			mockSetup:     func() {},
			expectedError: errs.ErrInvalidID,
		},
		{
			name: "missing user ID",
			ctx:  context.WithValue(context.Background(), domains.RoleKey{}, "buyer"),
			mockSetup: func() {
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
		{
			name: "get user error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			update: dto.UpdateUserProfileRequestDTO{
				Name: null.StringFrom("NewName"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(nil, errors.New("get user error"))
			},
			expectedError: errors.New("get user error"),
		},
		{
			name: "update error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUIDStr),
			update: dto.UpdateUserProfileRequestDTO{
				Name: null.StringFrom("NewName"),
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						Name: "OldName",
					}, nil)
				mockRepo.EXPECT().
					UpdateUserProfile(gomock.Any(), testUUID, gomock.Any()).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
		{
			name:          "invalid user ID",
			ctx:           context.WithValue(context.Background(), domains.UserIDKey{}, "invalid"),
			update:        dto.UpdateUserProfileRequestDTO{},
			mockSetup:     func() {},
			expectedError: errs.ErrInvalidID,
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

func TestUserUsecase_UpdateUserEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

	testUUID := uuid.New()
	testPassword := "password123"
	testEmail := "new@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		ctx           context.Context
		update        dto.UpdateUserEmailDTO
		mockSetup     func()
		expectedError error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserEmailDTO{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedPassword,
					}, nil)
				mockRepo.EXPECT().
					UpdateUserEmail(gomock.Any(), testUUID, testEmail).
					Return(nil)
			},
		},
		{
			name: "invalid password",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserEmailDTO{
				Email:    testEmail,
				Password: "wrongpassword",
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedPassword,
					}, nil)
			},
			expectedError: errs.ErrInvalidCredentials,
		},
		{
			name: "get user error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserEmailDTO{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(nil, errors.New("get user error"))
			},
			expectedError: errors.New("get user error"),
		},
		{
			name: "update error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserEmailDTO{
				Email:    testEmail,
				Password: testPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedPassword,
					}, nil)
				mockRepo.EXPECT().
					UpdateUserEmail(gomock.Any(), testUUID, testEmail).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := uc.UpdateUserEmail(tt.ctx, tt.update)
			if tt.expectedError != nil {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestUserUsecase_UpdateUserPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockMinio := minioMocks.NewMockProvider(ctrl)

	uc := user.NewUserUsecase(mockRepo, mockToken, mockMinio)

	testUUID := uuid.New()
	oldPassword := "oldPassword123"
	newPassword := "newPassword123"
	hashedOldPassword, _ := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		ctx           context.Context
		update        dto.UpdateUserPasswordDTO
		mockSetup     func()
		expectedError error
	}{
		{
			name: "success",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserPasswordDTO{
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedOldPassword,
					}, nil)
				mockRepo.EXPECT().
					UpdateUserPassword(gomock.Any(), testUUID, gomock.Any()).
					DoAndReturn(func(ctx context.Context, id uuid.UUID, hash []byte) error {
						// Verify the new password was hashed correctly
						err := bcrypt.CompareHashAndPassword(hash, []byte(newPassword))
						assert.NoError(t, err)
						return nil
					})
			},
		},
		{
			name: "invalid old password",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserPasswordDTO{
				OldPassword: "wrongpassword",
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedOldPassword,
					}, nil)
			},
			expectedError: errs.ErrInvalidCredentials,
		},
		{
			name: "get user error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserPasswordDTO{
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(nil, errors.New("get user error"))
			},
			expectedError: errors.New("get user error"),
		},
		{
			name: "update error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, testUUID.String()),
			update: dto.UpdateUserPasswordDTO{
				OldPassword: oldPassword,
				NewPassword: newPassword,
			},
			mockSetup: func() {
				mockRepo.EXPECT().
					GetUserByID(gomock.Any(), testUUID).
					Return(&models.UserDB{
						PasswordHash: hashedOldPassword,
					}, nil)
				mockRepo.EXPECT().
					UpdateUserPassword(gomock.Any(), testUUID, gomock.Any()).
					Return(errors.New("update error"))
			},
			expectedError: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := uc.UpdateUserPassword(tt.ctx, tt.update)
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
		{
			name:          "invalid user ID",
			ctx:           context.WithValue(context.Background(), domains.UserIDKey{}, "invalid"),
			request:       dto.UpdateRoleRequest{},
			mockSetup:     func() {},
			expectedError: errs.ErrInvalidID,
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