package user

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks
type IAuthUsecase interface {
	Register(context.Context, dto.UserRegisterRequestDTO) (string, error)
	Login(context.Context, dto.UserLoginRequestDTO) (string, error)
	Logout(context.Context) error
	GetMe(context.Context) (*models.User, error)
	UploadAvatar(context.Context, minio.FileData) (string, error)
}

type AuthService struct {
	AuthService  IAuthUsecase
	minioService minio.Provider
	config       config.Config
}

func NewAuthService(
	u IAuthUsecase,
	ms minio.Provider,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		AuthService:  u,
		minioService: ms,
		config:       *cfg,
	}
}

//	@Summary		Login user
//	@Description	Авторизация пользователя
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.UserLoginRequestDTO	true	"User credentials"
//	@success		200		{}			-							"No Content"
//	@Header			200		{string}	Set-Set						"Устанавливает JWT-токен в куки"
//	@Failure		400		{object}	dto.ErrorResponse			"Ошибка валидации"
//	@Failure		401		{object}	dto.ErrorResponse			"Неверные email или пароль"
//	@Failure		500		{object}	dto.ErrorResponse			"Внутренняя ошибка сервера"
//	@Router			/auth/login [post]

func (h *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	const op = "AuthService.Login"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var loginReq dto.UserLoginRequestDTO
	if err := request.ParseData(r, &loginReq); err != nil {
        logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	validator.SanitizeUserLoginRequest(&loginReq)

	if err := validator.ValidateLoginCreds(loginReq); err != nil {
        logger.WithError(err).WithField("email", loginReq.Email).Warn("validation failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	token, err := h.AuthService.Login(r.Context(), loginReq)
	if err != nil {
        logger.WithError(err).WithField("email", loginReq.Email).Error("login failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Set(w, token, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

//	@Summary		Register user
//	@Description	Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.UserRegisterRequestDTO	true	"Данные для регистрации"
//	@success		200		{}			-							"No Content"
//	@Header			200		{string}	Set-Set						"Устанавливает JWT-токен в куки"
//	@Failure		400		{object}	dto.ErrorResponse			"Некорректный запрос"
//	@Failure		409		{object}	dto.ErrorResponse			"Пользователь уже существует"
//	@Failure		500		{object}	dto.ErrorResponse			"Внутренняя ошибка сервера"
//	@Router			/auth/register [post]
func (h *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	const op = "AuthService.Register"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    var registerReq dto.UserRegisterRequestDTO
    if err := request.ParseData(r, &registerReq); err != nil {
        logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	validator.SanitizeUserRegistrationRequest(&registerReq)

	if err := validator.ValidateRegistrationCreds(registerReq); err != nil {
        logger.WithError(err).WithField("email", registerReq.Email).Warn("validation failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    token, err := h.AuthService.Register(r.Context(), registerReq)
    if err != nil {
        logger.WithError(err).WithField("email", registerReq.Email).Error("registration failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Set(w, token, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

//	@Summary		Logout user
//	@Description	Выход пользователя
//	@Tags			auth
//	@Security		TokenAuth
//	@Success		200	{}			"No Content"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не найден"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/auth/logout [post]
func (h *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthService.Logout"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := h.AuthService.Logout(r.Context()); err != nil {
        logger.WithError(err).Error("logout failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Unset(w, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

//	@Summary		Get user info
//	@Description	Получение информации о текущем пользователе
//	@Tags			users
//	@Security		TokenAuth
//	@Produce		json
//	@Success		200	{object}	models.User			"Информация о пользователе"
//	@Failure		400	{object}	dto.ErrorResponse	"Некорректный запрос"
//	@Failure		401	{object}	dto.ErrorResponse	"Пользователь не найден"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/users/me [get]
func (h *AuthService) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "AuthService.GetMe"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	user, err := h.AuthService.GetMe(r.Context())
    if err != nil {
        logger.WithError(err).Error("get user info failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

	response.SendJSONResponse(r.Context(), w, http.StatusOK, user)
}

//	@Summary		Upload avatar
//	@Description	Загружает аватар пользователя
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file				true	"Файл изображения"
//	@Success		200		{object}	map[string]string	"URL загруженного аватара"
//	@Failure		400		{object}	dto.ErrorResponse	"Ошибка при обработке формы"
//	@Failure		500		{object}	dto.ErrorResponse	"Ошибка загрузки файла"
//	@Security		TokenAuth
//	@Router			/users/avatar [post]
func (h *AuthService) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	const op = "AuthService.UploadAvatar"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := r.ParseMultipartForm(h.config.ServerConfig.MaxMultipartMemory); err != nil {
        logger.WithError(err).Error("parse multipart form")
        response.HandleDomainError(r.Context(), w, fmt.Errorf("parse form data"), op)
        return
    }

    file, header, err := r.FormFile(h.config.ServerConfig.AvatarKey)
    if err != nil {
        logger.WithError(err).Error("get file from form")
        response.HandleDomainError(r.Context(), w, fmt.Errorf("no file uploaded"), op)
        return
    }

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
        logger.WithError(err).Error("read file content")
        response.HandleDomainError(r.Context(), w, fmt.Errorf("read file"), op)
        return
    }

	fileData := minio.FileData{
		Name: header.Filename,
		Data:     fileBytes,
	}

	logger = logger.WithField("filename", header.Filename)
    avatarURL, err := h.AuthService.UploadAvatar(r.Context(), fileData)
    if err != nil {
        logger.WithError(err).Error("upload avatar failed")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    logger.Debug("avatar uploaded successfully")
    response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]string{
        "imageURL": avatarURL,
    })
}
