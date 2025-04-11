package auth

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/validator"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:generate mockgen -source=auth.go -destination=../../usecase/mocks/auth_usecase_mock.go -package=mocks
type IAuthUsecase interface {
	Register(context.Context, dto.UserRegisterRequestDTO) (string, error)
	Login(context.Context, dto.UserLoginRequestDTO) (string, error)
	Logout(context.Context) error
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
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: u,
		log:         log,
		config:      *cfg,
	}
}

//	@Summary		Login auth
//	@Description	Авторизация пользователя
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.UserLoginRequestDTO	true	"User credentials"
//	@success		200		{}			-							"No Content"
//	@Header			200		{string}	Set-Set						"Устанавливает JWT-токен в куки"
//	@Failure		400		{object}	response.ErrorResponse		"Ошибка валидации"
//	@Failure		401		{object}	response.ErrorResponse		"Неверные email или пароль"
//	@Failure		500		{object}	response.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/auth/login [post]

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

	cookieProvider.Set(w, token, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary		Register auth
// @Description	Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		models.UserRegisterRequestDTO	true	"Данные для регистрации"
// @success		200		{}			-								"No Content"
// @Header			200		{string}	Set-Set							"Устанавливает JWT-токен в куки"
// @Failure		400		{object}	response.ErrorResponse			"Некорректный запрос"
// @Failure		409		{object}	response.ErrorResponse			"Пользователь уже существует"
// @Failure		500		{object}	response.ErrorResponse			"Внутренняя ошибка сервера"
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

	cookieProvider.Set(w, token, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

// @Summary		Logout auth
// @Description	Выход пользователя
// @Tags			auth
// @Security		TokenAuth
// @Success		200	{}			"No Content"
// @Failure		401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.authService.Logout(r.Context()); err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to logout")
		return
	}

	cookieProvider := cookie.NewCookieProvider(&h.config)

	cookieProvider.Unset(w, string(domains.Token))

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
