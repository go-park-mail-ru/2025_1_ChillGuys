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
	UploadAvatar(context.Context, minio.FileData) (string, error)
	UpdateUserProfile(context.Context, dto.UpdateUserProfileRequestDTO) error
	UpdateUserEmail(ctx context.Context, user dto.UpdateUserEmailDTO) error
	UpdateUserPassword(context.Context, dto.UpdateUserPasswordDTO) error
}

type UserHandler struct {
	userService  IUserUsecase
	log          *logrus.Logger
	minioService minio.Provider
	config       config.Config
}

func NewUserHandler(
	u IUserUsecase,
	log *logrus.Logger,
	ms minio.Provider,
	cfg *config.Config,
) *UserHandler {
	return &UserHandler{
		userService:  u,
		log:          log,
		minioService: ms,
		config:       *cfg,
	}
}

//	@Summary		Получить информацию о себе
//	@Description	Возвращает информацию о текущем авторизованном пользователе
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	dto.UserDTO	"Информация о пользователе"
//	@Failure		400	{string}	string		"Некорректный запрос"
//	@Failure		401	{string}	string		"Пользователь не авторизован"
//	@Failure		500	{string}	string		"Внутренняя ошибка сервера"
//	@Router			/users/me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.userService.GetMe(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get current user")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

//	@Summary		Загрузить аватар
//	@Description	Загружает изображение профиля пользователя
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file				true	"Файл изображения"
//	@Success		200		{object}	map[string]string	"URL загруженного аватара"
//	@Failure		400		{string}	string				"Ошибка загрузки или обработки формы"
//	@Failure		500		{string}	string				"Внутренняя ошибка сервера"
//	@Router			/users/avatar [post]
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

	fileData := minio.FileData{
		Name: header.Filename,
		Data: fileBytes,
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

//	@Summary		Обновить профиль пользователя
//	@Description	Обновляет основную информацию пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateUserProfileRequestDTO	true	"Данные для обновления профиля"
//	@Success		200		{string}	string							"Профиль успешно обновлён"
//	@Failure		400		{string}	string							"Невалидные данные"
//	@Failure		500		{string}	string							"Ошибка при обновлении профиля"
//	@Security		TokenAuth
//	@Router			/users/update-profile [post]
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

//	@Summary		Обновить email пользователя
//	@Description	Обновляет email текущего пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateUserEmailDTO	true	"Новый email"
//	@Success		200		{string}	string					"Email успешно обновлён"
//	@Failure		400		{string}	string					"Невалидные данные"
//	@Failure		500		{string}	string					"Ошибка при обновлении email"
//	@Security		TokenAuth
//	@Router			/users/update-email [post]
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

//	@Summary		Обновить пароль пользователя
//	@Description	Меняет пароль текущего пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.UpdateUserPasswordDTO	true	"Старый и новый пароли"
//	@Success		200		{string}	string						"Пароль успешно обновлён"
//	@Failure		400		{string}	string						"Невалидные данные"
//	@Failure		500		{string}	string						"Ошибка при обновлении пароля"
//	@Security		TokenAuth
//	@Router			/users/update-password [post]
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
