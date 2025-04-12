package user

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks IUserUsecase
type IUserUsecase interface {
	GetMe(context.Context) (*dto.UserDTO, error)
	UploadAvatar(context.Context, minio.FileDataType) (string, error)
	UpdateUserProfile(context.Context, dto.UpdateUserProfileRequestDTO) error
	UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error
	UpdateUserPassword(context.Context, dto.UpdateUserPasswordDTO) error
}

type UserHandler struct {
	userService  IUserUsecase
	log          *logrus.Logger
	minioService minio.Client
	config       config.Config
}

func NewUserHandler(
	u IUserUsecase,
	log *logrus.Logger,
	ms minio.Client,
	cfg *config.Config,
) *UserHandler {
	return &UserHandler{
		userService:  u,
		log:          log,
		minioService: ms,
		config:       *cfg,
	}
}

// @Summary		Get auth info
// @Description	Получение информации о текущем пользователе
// @Tags			users
// @Security		TokenAuth
// @Produce		json
// @Success		200	{object}	models.UserDTO				"Информация о пользователе"
// @Failure		400	{object}	response.ErrorResponse	"Некорректный запрос"
// @Failure		401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetMe(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get current user")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

// @Summary		Upload avatar
// @Description	Загружает аватар пользователя
// @Tags			users
// @Accept			multipart/form-data
// @Produce		json
// @Param			file	formData	file					true	"Файл изображения"
// @Success		200		{object}	map[string]string		"URL загруженного аватара"
// @Failure		400		{object}	response.ErrorResponse	"Ошибка при обработке формы"
// @Failure		500		{object}	response.ErrorResponse	"Ошибка загрузки файла"
// @Security		TokenAuth
// @Router			/users/avatar [post]
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(h.config.ServerConfig.MaxMultipartMemory); err != nil {
		h.log.Warnf("error parsing multipart form: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, header, err := r.FormFile(h.config.ServerConfig.AvatarKey)
	if err != nil {
		h.log.Warnf("error getting file from form: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no file uploaded")
		return
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.log.Errorf("error reading file: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to read file")
		return
	}

	fileData := minio.FileDataType{
		FileName: header.Filename,
		Data:     fileBytes,
	}

	avatarURL, err := h.userService.UploadAvatar(r.Context(), fileData)
	if err != nil {
		h.log.Errorf("upload error: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "upload failed")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]string{
		"imageURL": avatarURL,
	})
}

func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var updateReq dto.UpdateUserProfileRequestDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserProfileUpdateRequest(&updateReq)

	if err := validator.ValidateUserUpdateProfileCreds(updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateUserProfile(r.Context(), updateReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to update user profile")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

func (h *UserHandler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var updateReq dto.UpdateUserEmailDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserEmailUpdateRequest(&updateReq)

	if err := validator.ValidateEmailCreds(updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateUserEmail(r.Context(), updateReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to update user email")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	var updateReq dto.UpdateUserPasswordDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserPasswordUpdateRequest(&updateReq)

	if err := validator.ValidatePasswordCreds(updateReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.userService.UpdateUserPassword(r.Context(), updateReq); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to update user password")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
