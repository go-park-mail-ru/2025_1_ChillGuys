package user

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks
type IAuthUsecase interface {
	Register(context.Context, dto.UserRegisterRequestDTO) (string, error)
	Login(context.Context, dto.UserLoginRequestDTO) (string, error)
	Logout(context.Context) error
	GetMe(context.Context) (*models.User, error)
	UploadAvatar(context.Context, minio.FileDataType) (string, error)
}

type AuthHandler struct {
	authService  IAuthUsecase
	log          *logrus.Logger
	minioService minio.Client
	config       config.Config
}

func NewAuthHandler(
	u IAuthUsecase,
	log *logrus.Logger,
	ms minio.Client,
) *AuthHandler {
	return &AuthHandler{
		authService:  u,
		log:          log,
		minioService: ms,
	}
}

// @Summary		Login user
// @Description	Авторизация пользователя
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		models.UserLoginRequestDTO	true	"User credentials"
// @success		200		{}			-							"No Content"
// @Header			200		{string}	Set-Set					"Устанавливает JWT-токен в куки"
// @Failure		400		{object}	response.ErrorResponse		"Ошибка валидации"
// @Failure		401		{object}	response.ErrorResponse		"Неверные email или пароль"
// @Failure		500		{object}	response.ErrorResponse		"Внутренняя ошибка сервера"
// @Router			/auth/login [post]

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq dto.UserLoginRequestDTO
	if err := request.ParseData(r, &loginReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserLoginRequest(&loginReq)

	if err := validator.ValidateLoginCreds(loginReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Login(r.Context(), loginReq)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to login user")
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Set(w, token, domains.Token.String())

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary			Register user
// @Description		Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
// @Tags			auth
// @Accept			json
// @Produce			json
// @Param			input	body		models.UserRegisterRequestDTO	true	"Данные для регистрации"
// @success			200		{}			-								"No Content"
// @Header			200		{string}	Set-Set						"Устанавливает JWT-токен в куки"
// @Failure			400		{object}	response.ErrorResponse			"Некорректный запрос"
// @Failure			409		{object}	response.ErrorResponse			"Пользователь уже существует"
// @Failure			500		{object}	response.ErrorResponse			"Внутренняя ошибка сервера"
// @Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerReq dto.UserRegisterRequestDTO
	if err := request.ParseData(r, &registerReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserRegistrationRequest(&registerReq)

	if err := validator.ValidateRegistrationCreds(registerReq); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	// Вызов usecase для регистрации
	token, err := h.authService.Register(r.Context(), registerReq)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to register user")
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Set(w, token, domains.Token.String())

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary			Logout user
// @Description		Выход пользователя
// @Tags			auth
// @Security		TokenAuth
// @Success			200	{}			"No Content"
// @Failure			401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure			500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.authService.Logout(r.Context()); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to logout")
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Unset(w, domains.Token.String())

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary			Get user info
// @Description		Получение информации о текущем пользователе
// @Tags			users
// @Security		TokenAuth
// @Produce			json
// @Success			200	{object}	models.User				"Информация о пользователе"
// @Failure			400	{object}	response.ErrorResponse	"Некорректный запрос"
// @Failure			401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure			500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.authService.GetMe(r.Context())
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get current user")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

// @Summary			Upload avatar
// @Description		Загружает аватар пользователя
// @Tags			users
// @Accept			multipart/form-data
// @Produce			json
// @Param			file	formData	file					true	"Файл изображения"
// @Success			200		{object}	map[string]string		"URL загруженного аватара"
// @Failure			400		{object}	response.ErrorResponse	"Ошибка при обработке формы"
// @Failure			500		{object}	response.ErrorResponse	"Ошибка загрузки файла"
// @Security		TokenAuth
// @Router			/users/avatar [post]
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
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

	avatarURL, err := h.authService.UploadAvatar(r.Context(), fileData)
	if err != nil {
		h.log.Errorf("upload error: %v", err)
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "upload failed")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, avatarURL)
}
