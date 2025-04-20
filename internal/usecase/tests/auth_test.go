package tests

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAuthRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockToken)

	tests := []struct {
		name          string
		input         dto.UserRegisterRequestDTO
		mockRepoSetup func()
		expectedToken string
		expectedID    uuid.UUID
		expectedErr   error
	}{
		{
			name: "Successful registration",
			input: dto.UserRegisterRequestDTO{
				Email:    "test@example.com",
				Password: "password",
				Name:     "Test",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().CheckUserExists(gomock.Any(), "test@example.com").Return(false, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedToken: "token",
			expectedID:    uuid.New(),
			expectedErr:   nil,
		},
		{
			name: "User already exists",
			input: dto.UserRegisterRequestDTO{
				Email:    "existing@example.com",
				Password: "password",
				Name:     "Existing",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().CheckUserExists(gomock.Any(), "existing@example.com").Return(true, nil)
			},
			expectedToken: "",
			expectedID:    uuid.Nil,
			expectedErr:   errors.New("AuthUsecase.Register: already exists"),
		},
		{
			name: "Repository error on check",
			input: dto.UserRegisterRequestDTO{
				Email:    "error@example.com",
				Password: "password",
				Name:     "Error",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().CheckUserExists(gomock.Any(), "error@example.com").Return(false, errors.New("repo error"))
			},
			expectedToken: "",
			expectedID:    uuid.Nil,
			expectedErr:   errors.New("AuthUsecase.Register: repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			if tt.expectedErr == nil {
				mockToken.EXPECT().CreateJWT(gomock.Any(), 1).Return(tt.expectedToken, nil)
			}

			token, id, err := authUC.Register(context.Background(), tt.input)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedID != uuid.Nil {
				assert.NotEqual(t, uuid.Nil, id)
			} else {
				assert.Equal(t, tt.expectedID, id)
			}

			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAuthRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockToken)

	testUserID := uuid.New()
	testUser := &models.UserDB{
		ID:    testUserID,
		Email: "test@example.com",
		PasswordHash: func() []byte {
			hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
			return hash
		}(),
		UserVersion: models.UserVersionDB{
			Version: 1,
		},
	}

	tests := []struct {
		name          string
		input         dto.UserLoginRequestDTO
		mockRepoSetup func()
		expectedToken string
		expectedID    uuid.UUID
		expectedErr   error
	}{
		{
			name: "Successful login",
			input: dto.UserLoginRequestDTO{
				Email:    "test@example.com",
				Password: "password",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(testUser, nil)
			},
			expectedToken: "token",
			expectedID:    testUserID,
			expectedErr:   nil,
		},
		{
			name: "User not found",
			input: dto.UserLoginRequestDTO{
				Email:    "notfound@example.com",
				Password: "password",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "notfound@example.com").Return(nil, errs.ErrNotFound)
			},
			expectedToken: "",
			expectedID:    uuid.Nil,
			expectedErr:   errs.ErrNotFound,
		},
		{
			name: "Invalid password",
			input: dto.UserLoginRequestDTO{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "test@example.com").Return(testUser, nil)
			},
			expectedToken: "",
			expectedID:    uuid.Nil,
			expectedErr:   errs.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			if tt.expectedErr == nil {
				mockToken.EXPECT().CreateJWT(testUserID.String(), 1).Return(tt.expectedToken, nil)
			}

			token, id, err := authUC.Login(context.Background(), tt.input)

			assert.Equal(t, tt.expectedToken, token)
			assert.Equal(t, tt.expectedID, id)

			if tt.expectedErr != nil {
				assert.True(t, errors.Is(err, tt.expectedErr), "expected error %v, got %v", tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAuthRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockToken)

	tests := []struct {
		name          string
		ctx           context.Context
		mockRepoSetup func()
		expectedErr   error
	}{
		{
			name: "Successful logout",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, "user-id"),
			mockRepoSetup: func() {
				mockRepo.EXPECT().IncrementUserVersion(gomock.Any(), "user-id").Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:          "User ID not in context",
			ctx:           context.Background(),
			mockRepoSetup: func() {},
			expectedErr:   errs.ErrNotFound,
		},
		{
			name: "Repository error",
			ctx:  context.WithValue(context.Background(), domains.UserIDKey{}, "user-id"),
			mockRepoSetup: func() {
				mockRepo.EXPECT().IncrementUserVersion(gomock.Any(), "user-id").Return(errors.New("repo error"))
			},
			expectedErr: errors.New("repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			err := authUC.Logout(tt.ctx)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				// Check either exact match or that the error contains our expected message
				assert.True(t, errors.Is(err, tt.expectedErr) ||
					strings.Contains(err.Error(), tt.expectedErr.Error()),
					"expected error containing %q, got %v", tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGeneratePasswordHash(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "Valid password",
			password:    "validpassword",
			expectedErr: nil,
		},
		{
			name:        "Empty password",
			password:    "",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := auth.GeneratePasswordHash(tt.password)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
