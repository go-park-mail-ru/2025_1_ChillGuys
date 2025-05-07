package user

import (
	"context"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/user"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/guregu/null"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGRPCHandler struct {
	gen.UnimplementedUserServiceServer
	
	userProvider user.IUserUsecase
	minioService minio.Provider
}

func NewUserGRPCHandler(u user.IUserUsecase, ms minio.Provider) *UserGRPCHandler {
	return &UserGRPCHandler{
		userProvider: u,
		minioService: ms,
	}
}

func (h *UserGRPCHandler) GetMe(ctx context.Context, _ *emptypb.Empty) (*gen.User, error){
	const op = "UserGRPCHandler.GetMe"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	user, token, err := h.userProvider.GetMe(ctx)
	if err != nil {
		logger.WithError(err).Error("get current user")
		return nil, errs.MapErrorToGRPC(err)
	}

	if token != "" {
        if err := grpc.SetHeader(ctx, metadata.Pairs("x-new-token", token)); err != nil {
            logger.WithError(err).Error("failed to set new token header")
        }
    }

	return user.ConvertToGrpcUser(), nil
}

func (h *UserGRPCHandler) UploadAvatar(stream gen.UserService_UploadAvatarServer) error {
	const op = "UserGRPCHandler.UploadAvatar"
	ctx := stream.Context()
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var fileData []byte
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.WithError(err).Error("failed to receive chunk")
			return errs.MapErrorToGRPC(errs.ErrInternal)
		}
		fileData = append(fileData, chunk.Value...)
	}

	contentType := http.DetectContentType(fileData)
	if err := validator.ValidateImageContentType(contentType); err != nil {
		logger.WithField("contentType", contentType).Error(err.Error())
		return status.Error(codes.InvalidArgument, err.Error())
	}

	avatarURL, err := h.userProvider.UploadAvatar(ctx, minio.FileData{
		Name: "avatar",
		Data: fileData,
	})
	if err != nil {
		logger.WithError(err).Error("upload failed")
		return errs.MapErrorToGRPC(err)
	}

	return stream.SendAndClose(&gen.UploadAvatarResponse{
		ImageURL: avatarURL,
	})
}

func (h *UserGRPCHandler) UpdateUserProfile(ctx context.Context, req *gen.UpdateUserProfileRequest) (*emptypb.Empty, error) {
	const op = "UserGRPCHandler.UpdateUserProfile"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	var name null.String
	if req.Name != nil {
		name = null.StringFrom(req.Name.Value)
	}

	var surname null.String
	if req.Surname != nil {
		surname = null.StringFrom(req.Surname.Value)
	}

	var phoneNumber null.String
	if req.PhoneNumber != nil {
		phoneNumber = null.StringFrom(req.PhoneNumber.Value)
	}

	request := dto.UpdateUserProfileRequestDTO{
		Name:        name,
		Surname:     surname,
		PhoneNumber: phoneNumber,
	}

	validator.SanitizeUserProfileUpdateRequest(&request)
	if err := validator.ValidateUserUpdateProfileCreds(request); err != nil {
		logger.WithError(err).Error("validate profile update credentials")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.userProvider.UpdateUserProfile(ctx, request); err != nil {
		logger.WithError(err).Error("update profile failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) UpdateUserEmail(ctx context.Context, req *gen.UpdateUserEmailRequest) (*emptypb.Empty, error) {
	const op = "UserGRPCHandler.UpdateUserEmail"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	request := dto.UpdateUserEmailDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	validator.SanitizeUserEmailUpdateRequest(&request)
	if err := validator.ValidateEmailCreds(request); err != nil {
		logger.WithError(err).Error("validate email update credentials")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.userProvider.UpdateUserEmail(ctx, request); err != nil {
		logger.WithError(err).Error("update email failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) UpdateUserPassword(ctx context.Context, req *gen.UpdateUserPasswordRequest) (*emptypb.Empty, error) {
	const op = "UserGRPCHandler.UpdateUserPassword"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	request := dto.UpdateUserPasswordDTO{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}

	validator.SanitizeUserPasswordUpdateRequest(&request)
	if err := validator.ValidatePasswordCreds(request); err != nil {
		logger.WithError(err).Error("validate password update credentials")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.userProvider.UpdateUserPassword(ctx, request); err != nil {
		logger.WithError(err).Error("update password failed")
		return nil, errs.MapErrorToGRPC(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserGRPCHandler) BecomeSeller(ctx context.Context, req *gen.BecomeSellerRequest) (*emptypb.Empty, error) {
    const op = "UserGRPCHandler.BecomeSeller"
    logger := logctx.GetLogger(ctx).WithField("op", op)

    request := dto.UpdateRoleRequest{
        Title:       req.Title,
        Description: req.Description,
    }

    if err := h.userProvider.BecomeSeller(ctx, request); err != nil {
        logger.WithError(err).Error("become seller failed")
        return nil, errs.MapErrorToGRPC(err)
    }

    return &emptypb.Empty{}, nil
}