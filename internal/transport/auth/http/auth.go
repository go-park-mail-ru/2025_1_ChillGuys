package http

import (
	"github.com/mailru/easyjson"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/auth"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/metadata"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthHandler struct {
	authClient gen.AuthServiceClient
	config     *config.Config
}

func NewAuthHandler(
	authClient gen.AuthServiceClient,
	cfg *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
		config:     cfg,
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
	if err := easyjson.UnmarshalFromReader(r.Body, &loginReq); err != nil {
		logger.WithError(err).Error("parse login request")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.authClient.Login(r.Context(), &gen.LoginReq{
		Email:    loginReq.Email,
		Password: loginReq.Password,
	})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken(
		res.Token,
		h.config.CSRFConfig.SecretKey,
		h.config.CSRFConfig.TokenExpiry,
	)
	if err != nil {
		logger.WithError(err).Error("generate CSRF token")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to generate CSRF token")
		return
	}

	cookieProvider := cookie.NewCookieProvider(h.config)
	cookieProvider.Set(w, res.Token, domains.TokenCookieName)
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
	if err := easyjson.UnmarshalFromReader(r.Body, &registerReq); err != nil {
		logger.WithError(err).Error("parse register request")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.authClient.Register(r.Context(), registerReq.ConvertToGrpcRegisterReq())
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	csrfToken, err := middleware.GenerateCSRFToken(
		res.Token,
		h.config.CSRFConfig.SecretKey,
		h.config.CSRFConfig.TokenExpiry,
	)
	if err != nil {
		logger.WithError(err).Error("generate CSRF token")
		response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "failed to generate CSRF token")
		return
	}

	cookieProvider := cookie.NewCookieProvider(h.config)
	cookieProvider.Set(w, res.Token, domains.TokenCookieName)
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

	jwtCookie, err := r.Cookie(string(domains.TokenCookieName))
	if err != nil {
		response.SendJSONError(r.Context(), w, http.StatusUnauthorized, "JWT token required")
		return
	}

	ctxWithToken := metadata.InjectJWTIntoContext(r.Context(), jwtCookie.Value)
	_, err = h.authClient.Logout(ctxWithToken, &emptypb.Empty{})
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		return
	}

	cookieProvider := cookie.NewCookieProvider(h.config)
	cookieProvider.Unset(w, domains.TokenCookieName)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}
