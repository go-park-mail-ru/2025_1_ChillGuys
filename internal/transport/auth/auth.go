package auth

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/google/uuid"
)

//go:generate mockgen -source=auth.go -destination=../../usecase/mocks/auth_usecase_mock.go -package=mocks IAuthUsecase
type IAuthUsecase interface {
	Register(context.Context, dto.UserRegisterRequestDTO) (string, uuid.UUID, error)
	Login(context.Context, dto.UserLoginRequestDTO) (string, uuid.UUID, error)
	Logout(context.Context) error
}

type AuthHandler struct {
	authService  IAuthUsecase
	minioService minio.Provider
	config       config.Config
}

func NewAuthHandler(
	u IAuthUsecase,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: u,
		config:      *cfg,
	}
}

// Login godoc
//
//	@Summary		Авторизация пользователя
//	@Description	Авторизует пользователя и устанавливает JWT-токен в cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.UserLoginRequestDTO	true	"Данные для входа"
//	@Success		200		{}			-						"Успешная авторизация"
//	@Header			200		{string}	Set-Cookie				"JWT-токен авторизации"
//	@Failure		400		{object}	object					"Ошибка валидации данных"
//	@Failure		401		{object}	object					"Неверные email или пароль"
//	@Failure		500		{object}	object					"Внутренняя ошибка сервера"
//	@Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Login"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)
	
	var loginReq dto.UserLoginRequestDTO
	if err := request.ParseData(r, &loginReq); err != nil {
		logger.WithError(err).Error("parse login request")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserLoginRequest(&loginReq)

	if err := validator.ValidateLoginCreds(loginReq); err != nil {
		logger.WithError(err).Error("validate login credentials")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	token, userID, err := h.authService.Login(r.Context(), loginReq)
	if err != nil {
		logger.WithError(err).Error("login failed")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	// Генерируем CSRF токен
	csrfToken, err := middleware.GenerateCSRFToken(
		token, // Используем JWT токен как sessionID
		userID,
		h.config.CSRFConfig.SecretKey,
		h.config.CSRFConfig.TokenExpiry,
	)
	if err != nil {
		logger.WithError(err).Error("generate CSRF token")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to generate CSRF token")
		return
	}

	// Устанавливаем JWT в куки
	cookieProvider := cookie.NewCookieProvider(&h.config)
	cookieProvider.Set(w, token, domains.TokenCookieName)

	// Устанавливаем CSRF токен в заголовок
	w.Header().Set("X-CSRF-Token", csrfToken)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// Register godoc
//
//	@Summary		Регистрация пользователя
//	@Description	Создает нового пользователя и устанавливает JWT-токен в cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			userData	body		dto.UserRegisterRequestDTO	true	"Данные для регистрации"
//	@Success		200			{}			-							"Успешная регистрация"
//	@Header			200			{string}	Set-Cookie					"JWT-токен авторизации"
//	@Failure		400			{object}	object						"Некорректные данные"
//	@Failure		409			{object}	object						"Пользователь уже существует"
//	@Failure		500			{object}	object						"Внутренняя ошибка сервера"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Register"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)
	
	var registerReq dto.UserRegisterRequestDTO
	if err := request.ParseData(r, &registerReq); err != nil {
		logger.WithError(err).Error("parse register request")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	validator.SanitizeUserRegistrationRequest(&registerReq)

	if err := validator.ValidateRegistrationCreds(registerReq); err != nil {
		logger.WithError(err).Error("validate registration credentials")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	token, userID, err := h.authService.Register(r.Context(), registerReq)
	if err != nil {
		logger.WithError(err).Error("register user failed")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	// Генерируем CSRF токен
	csrfToken, err := middleware.GenerateCSRFToken(
		token, // Используем JWT токен как sessionID
		userID,
		h.config.CSRFConfig.SecretKey,
		h.config.CSRFConfig.TokenExpiry,
	)
	if err != nil {
		logger.WithError(err).Error("generate CSRF token")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to generate CSRF token")
		return
	}

	// Устанавливаем JWT в куки
	cookieProvider := cookie.NewCookieProvider(&h.config)
	cookieProvider.Set(w, token, domains.TokenCookieName)

	// Устанавливаем CSRF токен в заголовок
	w.Header().Set("X-CSRF-Token", csrfToken)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// Logout godoc
//
//	@Summary		Выход из системы
//	@Description	Завершает сеанс пользователя и удаляет JWT-токен из cookies
//	@Tags			auth
//	@Produce		json
//	@Param			X-Csrf-Token	header		string		true	"CSRF-токен для защиты от подделки запросов"
//	@Success		200				{}			-			"Успешный выход из системы"
//	@Header			200				{string}	Set-Cookie	"Очищает JWT-токен (устанавливает пустое значение с истекшим сроком)"
//	@Failure		401				{object}	object		"Пользователь не авторизован"
//	@Failure		500				{object}	object		"Внутренняя ошибка сервера"
//	@Security		TokenAuth
//	@Router			/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "AuthHandler.Logout"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	if err := h.authService.Logout(r.Context()); err != nil {
		logger.WithError(err).Error("logout failed")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Unset(w, domains.TokenCookieName)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
