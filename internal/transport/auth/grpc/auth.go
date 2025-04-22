package auth

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/repository/redis"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/metadata"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthGRPCHandler struct {
	gen.UnimplementedAuthServiceServer

	authProvider auth.IAuthUsecase
	redisRepo    *redis.AuthRepository
	tokenator    *jwt.Tokenator
}

func NewAuthGRPCHandler(u auth.IAuthUsecase, redisRepo *redis.AuthRepository, tokenator *jwt.Tokenator) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		authProvider: u,
		redisRepo:    redisRepo,
		tokenator:    tokenator,
	}
}

func (h *AuthGRPCHandler) Register(ctx context.Context, in *gen.RegisterReq) (*gen.RegisterRes, error) {
	const op = "AuthGRPCHandler.Register"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	var surname null.String
	if in.Surname != nil {
		surname = null.StringFrom(in.Surname.Value)
	}

	request := dto.UserRegisterRequestDTO{
		Email:    in.Email,
		Password: in.Password,
		Name:     in.Name,
		Surname:  surname,
	}

	// Валидация
	validator.SanitizeUserRegistrationRequest(&request)
	if err := validator.ValidateRegistrationCreds(request); err != nil {
		logger.WithError(err).Error("validate registration credentials")
		// FIXME: Сделать маппинг ошибок
		return nil, err
	}

	token, err := h.authProvider.Register(ctx, request)
	if err != nil {
		// FIXME: Сделать маппинг ошибок
		return nil, err
	}

	return &gen.RegisterRes{
		Token: token,
	}, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, in *gen.LoginReq) (*gen.LoginRes, error) {
	const op = "AuthGRPCHandler.Login"
	logger := logctx.GetLogger(ctx).WithField("op", op)
	request := dto.UserLoginRequestDTO{
		Email:    in.Email,
		Password: in.Password,
	}

	// Валидация
	validator.SanitizeUserLoginRequest(&request)
	if err := validator.ValidateLoginCreds(request); err != nil {
		logger.WithError(err).Error("validate registration credentials")
		// FIXME: Сделать маппинг ошибок
		return nil, err
	}

	token, err := h.authProvider.Login(ctx, request)
	if err != nil {
		// FIXME: Сделать маппинг ошибок
		return nil, err
	}

	return &gen.LoginRes{
		Token: token,
	}, nil
}

func (h *AuthGRPCHandler) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "AuthGRPCHandler.Logout"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	token, err := metadata.ExtractJWTFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to extract token")
		return nil, errs.ErrInternal
	}

	claims, err := h.tokenator.ParseJWT(token)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		return nil, errs.ErrInvalidToken
	}

	// Добавляем токен в черный список с привязкой к userID
	if err := h.redisRepo.AddToBlacklist(ctx, claims.UserID, token); err != nil {
		logger.WithError(err).Error("failed to add token to blacklist")
		return nil, errs.ErrInternal
	}

	return &emptypb.Empty{}, nil
}
