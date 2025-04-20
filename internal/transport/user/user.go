package user

import (
	"context"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks IUserUsecase
type IUserUsecase interface {
	GetMe(context.Context) (*dto.UserDTO, error)
	UploadAvatar(context.Context, minio.FileData) (string, error)
	UpdateUserProfile(context.Context, dto.UpdateUserProfileRequestDTO) error
	UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error
	UpdateUserPassword(context.Context, dto.UpdateUserPasswordDTO) error
}

type UserHandler struct {
	userService  IUserUsecase
	minioService minio.Provider
	config       config.Config
}

func NewUserHandler(
	u IUserUsecase,
	ms minio.Provider,
	cfg *config.Config,
) *UserHandler {
	return &UserHandler{
		userService:  u,
		minioService: ms,
		config:       *cfg,
	}
}

// @Summary		Получить информацию о себе
// @Description	Возвращает информацию о текущем авторизованном пользователе
// @Tags			users
// @Produce		json
// @Success		200	{object}	dto.UserDTO	"Информация о пользователе"
// @Failure		400	{string}	string		"Некорректный запрос"
// @Failure		401	{string}	string		"Пользователь не авторизован"
// @Failure		500	{string}	string		"Внутренняя ошибка сервера"
// @Router			/users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.GetMe"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	user, err := h.userService.GetMe(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get current user")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

// @Summary		Загрузить аватар
// @Description	Загружает изображение профиля пользователя
// @Tags			users
// @Accept			multipart/form-data
// @Produce		json
// @Param			file			formData	file				true	"Файл изображения"
// @Param			X-Csrf-Token	header		string				true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{object}	map[string]string	"URL загруженного аватара"
// @Failure		400				{string}	string				"Ошибка загрузки или обработки формы"
// @Failure		500				{string}	string				"Внутренняя ошибка сервера"
// @Router			/users/avatar [post]
func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UploadAvatar"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := r.ParseMultipartForm(h.config.ServerConfig.MaxMultipartMemory); err != nil {
		logger.WithError(err).Error("failed to parse form data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, header, err := r.FormFile(h.config.ServerConfig.AvatarKey)
	if err != nil {
		logger.WithError(err).Error("no file uploaded")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("failed to read file")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to read file")
		return
	}

	contentType := http.DetectContentType(fileBytes)
	if err := validator.ValidateImageContentType(contentType); err != nil {
		logger.WithField("contentType", contentType).Error(err.Error())
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	fileData := minio.FileData{
		Name: header.Filename,
		Data: fileBytes,
	}

	avatarURL, err := h.userService.UploadAvatar(r.Context(), fileData)
	if err != nil {
		logger.WithError(err).Error("upload failed")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "upload failed")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]string{
		"imageURL": avatarURL,
	})
}

// @Summary		Обновить профиль пользователя
// @Description	Обновляет основную информацию пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserProfileRequestDTO	true	"Данные для обновления профиля"
// @Param			X-Csrf-Token	header		string							true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string							"Профиль успешно обновлён"
// @Failure		400				{string}	string							"Невалидные данные"
// @Failure		500				{string}	string							"Ошибка при обновлении профиля"
// @Security		TokenAuth
// @Router			/users/update-profile [post]
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserProfile"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserProfileRequestDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
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

// @Summary		Обновить email пользователя
// @Description	Обновляет email текущего пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserEmailDTO	true	"Новый email"
// @Param			X-Csrf-Token	header		string					true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string					"Email успешно обновлён"
// @Failure		400				{string}	string					"Невалидные данные"
// @Failure		500				{string}	string					"Ошибка при обновлении email"
// @Security		TokenAuth
// @Router			/users/update-email [post]
func (h *UserHandler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserEmail"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserEmailDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.HandleDomainError(r.Context(), w, err, op)
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

// @Summary		Обновить пароль пользователя
// @Description	Меняет пароль текущего пользователя
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			body			body		dto.UpdateUserPasswordDTO	true	"Старый и новый пароли"
// @Param			X-Csrf-Token	header		string						true	"CSRF-токен для защиты от подделки запросов"
// @Success		200				{string}	string						"Пароль успешно обновлён"
// @Failure		400				{string}	string						"Невалидные данные"
// @Failure		500				{string}	string						"Ошибка при обновлении пароля"
// @Security		TokenAuth
// @Router			/users/update-password [post]
func (h *UserHandler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.UpdateUserPassword"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var updateReq dto.UpdateUserPasswordDTO
	if err := request.ParseData(r, &updateReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.HandleDomainError(r.Context(), w, err, op)
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
