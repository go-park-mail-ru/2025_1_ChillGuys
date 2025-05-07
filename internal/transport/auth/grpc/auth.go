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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	validator.SanitizeUserRegistrationRequest(&request)
	if err := validator.ValidateRegistrationCreds(request); err != nil {
		logger.WithError(err).Error("validate registration credentials")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := h.authProvider.Register(ctx, request)
	if err != nil {
		logger.WithError(err).Error("registration failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &gen.RegisterRes{Token: token}, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, in *gen.LoginReq) (*gen.LoginRes, error) {
	const op = "AuthGRPCHandler.Login"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	request := dto.UserLoginRequestDTO{
		Email:    in.Email,
		Password: in.Password,
	}

	validator.SanitizeUserLoginRequest(&request)
	if err := validator.ValidateLoginCreds(request); err != nil {
		logger.WithError(err).Error("validate login credentials")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := h.authProvider.Login(ctx, request)
	if err != nil {
		logger.WithError(err).Error("login failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &gen.LoginRes{Token: token}, nil
}

func (h *AuthGRPCHandler) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	const op = "AuthGRPCHandler.Logout"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	token, err := metadata.ExtractJWTFromContext(ctx)
	if err != nil {
		logger.WithError(err).Error("failed to extract token")
		return nil, errs.MapErrorToGRPC(errs.ErrInvalidToken)
	}

	if err := h.authProvider.Logout(ctx, token); err != nil {
		logger.WithError(err).Error("logout failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthGRPCHandler) CheckToken(ctx context.Context, req *gen.CheckTokenReq) (*gen.CheckTokenRes, error) {
	const op = "AuthGRPCHandler.CheckToken"
	tokenString := req.Token
	if tokenString == "" {
		return nil, errs.MapErrorToGRPC(errs.ErrInvalidToken)
	}

	claims, err := h.tokenator.ParseJWT(tokenString)
	if err != nil {
		return nil, errs.MapErrorToGRPC(errs.ErrInvalidToken)
	}

	isInBlackList, err := h.redisRepo.IsInBlacklist(ctx, claims.UserID, tokenString)
	if err != nil {
		return nil, errs.MapErrorToGRPC(errs.ErrInternal)
	}

	if isInBlackList {
		return nil, errs.MapErrorToGRPC(errs.ErrTokenRevoked)
	}

	return &gen.CheckTokenRes{Valid: true}, nil
}
