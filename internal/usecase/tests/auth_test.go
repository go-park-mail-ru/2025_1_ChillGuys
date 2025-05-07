package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/postgres/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/auth"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIAuthRepository(ctrl)
	mockToken := mocks.NewMockITokenator(ctrl)
	mockRedis := mocks.NewMockIAuthRedisRepository(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockRedis, mockToken)

	tests := []struct {
		name          string
		input         dto.UserRegisterRequestDTO
		mockRepoSetup func()
		expectedToken string
		expectedErr   error
	}{
		{
			name: "Successful registration",
			input: dto.UserRegisterRequestDTO{
				Email:    "test@example.com",
				Password: "password",
				Name:     "Test",
				Surname:  null.StringFrom("User"),
			},
			mockRepoSetup: func() {
				mockRepo.EXPECT().CheckUserExists(gomock.Any(), "test@example.com").Return(false, nil)
				mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, user models.UserDB) error {
						// Verify password is hashed
						err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte("password"))
						assert.NoError(t, err)
						return nil
					})
				mockToken.EXPECT().CreateJWT(gomock.Any(), models.RoleBuyer.String()).Return("token", nil)
			},
			expectedToken: "token",
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
			expectedErr:   errs.ErrAlreadyExists,
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
			expectedErr:   errors.New("AuthUsecase.Register: repo error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			token, err := authUC.Register(context.Background(), tt.input)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedErr != nil {
				assert.ErrorContains(t, err, tt.expectedErr.Error())
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
	mockRedis := mocks.NewMockIAuthRedisRepository(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockRedis, mockToken)

	testUserID := uuid.New()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	testUser := &models.UserDB{
		ID:           testUserID,
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         models.RoleBuyer,
	}

	tests := []struct {
		name          string
		input         dto.UserLoginRequestDTO
		mockRepoSetup func()
		expectedToken string
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
				mockToken.EXPECT().CreateJWT(testUserID.String(), models.RoleBuyer.String()).Return("token", nil)
			},
			expectedToken: "token",
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
			expectedErr:   errs.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockRepoSetup()

			token, err := authUC.Login(context.Background(), tt.input)

			assert.Equal(t, tt.expectedToken, token)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
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
	mockRedis := mocks.NewMockIAuthRedisRepository(ctrl)

	authUC := auth.NewAuthUsecase(mockRepo, mockRedis, mockToken)

	testToken := "test_token"
	testClaims := &jwt.JWTClaims{
		UserID: "user_id",
	}

	tests := []struct {
		name          string
		token         string
		mockSetup     func()
		expectedErr   error
	}{
		{
			name:  "Successful logout",
			token: testToken,
			mockSetup: func() {
				mockToken.EXPECT().ParseJWT(testToken).Return(testClaims, nil)
				mockRedis.EXPECT().AddToBlacklist(gomock.Any(), testClaims.UserID, testToken).Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:  "Invalid token",
			token: "invalid_token",
			mockSetup: func() {
				mockToken.EXPECT().ParseJWT("invalid_token").Return(nil, errors.New("invalid token"))
			},
			expectedErr: errs.ErrInvalidToken,
		},
		{
			name:  "Redis error",
			token: testToken,
			mockSetup: func() {
				mockToken.EXPECT().ParseJWT(testToken).Return(testClaims, nil)
				mockRedis.EXPECT().AddToBlacklist(gomock.Any(), testClaims.UserID, testToken).Return(errors.New("redis error"))
			},
			expectedErr: errs.ErrInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := authUC.Logout(context.Background(), tt.token)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}