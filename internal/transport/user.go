package transport

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
)

var (
	emailRegexp     = regexp.MustCompile(`^\w+(\.\w*)*@\w+(\.\w{2,})+$`)
	nameRegexp      = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
	digitRegexp     = regexp.MustCompile(`[0-9]`)
	lowercaseRegexp = regexp.MustCompile(`[a-z]`)
	uppercaseRegexp = regexp.MustCompile(`[A-Z]`)
)

//go:generate mockgen -source=user.go -destination=../repository/mocks/user_repo_mock.go -package=mocks IUserRepository
type IUserRepository interface {
	CreateUser(user models.UserDB) error
	GetUserByEmail(email string) (*models.UserDB, error)
	GetUserByID(id uuid.UUID) (*models.UserDB, error)
	IncrementUserVersion(userID string) error
}

type ITokenator interface {
	CreateJWT(userID string, version int) (string, error)
	ParseJWT(tokenString string) (*jwt.JWTClaims, error)
}

type AuthHandler struct {
	repo  IUserRepository
	token ITokenator
	log   *logrus.Logger
}

func NewAuthHandler(repo IUserRepository, log *logrus.Logger, token ITokenator) *AuthHandler {
	return &AuthHandler{
		repo:  repo,
		token: token,
		log:   log,
	}
}

// @Summary			Login user
// @Description		Авторизация пользователя
// @Tags			auth
// @Accept			json
// @Produce			json
// @Param			request		body		models.UserLoginRequestDTO	true	"User credentials"
// @success			200			{}			-							"No Content"
// @Header			200			{string}	Set-Cookie					"Устанавливает JWT-токен в куки"
// @Failure			400			{object}	utils.ErrorResponse			"Ошибка валидации"
// @Failure			401			{object}	utils.ErrorResponse			"Неверные email или пароль"
// @Failure			500			{object}	utils.ErrorResponse			"Внутренняя ошибка сервера"
// @Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request models.UserLoginRequestDTO
	if errStatusCode, err := utils.ParseData(r.Body, &request); err != nil {
		utils.SendErrorResponse(w, errStatusCode, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	if err := ValidateLoginCreds(request); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userRepo, err := h.repo.GetUserByEmail(request.Email)
	if err != nil {
		h.log.Warn(err.Error())
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(request.Password)); err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	token, err := h.token.CreateJWT(userRepo.ID.String(), userRepo.Version)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Cookie(w, token, string(utils.Token))
	utils.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary			Register user
// @Description		Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
// @Tags			auth
// @Accept			json
// @Produce			json
// @Param			input	body		models.UserRegisterRequestDTO	true	"Данные для регистрации"
// @success			200		{}			-								"No Content"
// @Header			200		{string}	Set-Cookie						"Устанавливает JWT-токен в куки"
// @Failure			400		{object}	utils.ErrorResponse				"Некорректный запрос"
// @Failure			409		{object}	utils.ErrorResponse				"Пользователь уже существует"
// @Failure			500		{object}	utils.ErrorResponse				"Внутренняя ошибка сервера"
// @Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request models.UserRegisterRequestDTO
	if errStatusCode, err := utils.ParseData(r.Body, &request); err != nil {
		utils.SendErrorResponse(w, errStatusCode, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	if err := ValidateRegistrationCreds(request); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	passwordHash, err := GeneratePasswordHash(request.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	existedUser, _ := h.repo.GetUserByEmail(request.Email)
	if existedUser != nil {
		utils.SendErrorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	userRepo := models.UserDB{
		ID:           uuid.New(),
		Email:        request.Email,
		Name:         request.Name,
		Surname:      request.Surname,
		PasswordHash: passwordHash,
		Version:      1,
	}

	if err := h.repo.CreateUser(userRepo); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	token, err := h.token.CreateJWT(userRepo.ID.String(), 1)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Cookie(w, token, string(utils.Token))
	utils.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary			Logout user
// @Description		Выход пользователя
// @Tags			auth
// @Security		TokenAuth
// @Success			200	{}			"No Content"
// @Failure			500	{object}	utils.ErrorResponse	"Ошибка сервера"
// @Router			/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, isExist := r.Context().Value(utils.UserIDKey).(string)
	if !isExist {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "User id not found")
		return
	}

	if err := h.repo.IncrementUserVersion(userID); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     string(utils.Token),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().UTC().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   true,
	})

	utils.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary			Get user info
// @Description		Получение информации о текущем пользователе
// @Tags			users
// @Security		TokenAuth
// @Produce			json
// @Success			200	{object}	models.User			"Информация о пользователе"
// @Failure			400	{object}	utils.ErrorResponse	"Некорректный запрос"
// @Failure			404	{object}	utils.ErrorResponse	"Пользователь не найден"
// @Failure			500	{object}	utils.ErrorResponse	"Ошибка сервера"
// @Router			/users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDStr, isExist := r.Context().Value(utils.UserIDKey).(string)
	if !isExist {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "User id not found")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid user id format")
		return
	}

	userRepo, err := h.repo.GetUserByID(userID)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	user := userRepo.ConvertToUser()
	utils.SendSuccessResponse(w, http.StatusOK, user)
}

// ValidateLoginCreds проверяет корректность данных при авторизации
func ValidateLoginCreds(req models.UserLoginRequestDTO) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// ValidateRegistrationCreds проверяет корректность данных при регистрации
func ValidateRegistrationCreds(req models.UserRegisterRequestDTO) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if err := validateName(req.Name); err != nil {
		return err
	}

	if req.Surname.Valid && strings.TrimSpace(req.Surname.String) != "" {
		if err := validateName(req.Surname.String); err != nil {
			return err
		}
	}

	return nil
}

// validateEmail Функция валидации почты
func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("Invalid email")
	}
	return nil
}

// validatePassword проверяет валидность пароля
func validatePassword(password string) error {
	switch {
	case len(password) < 8:
		return errors.New("Password must be at least 8 characters")
	case !digitRegexp.MatchString(password):
		return errors.New("Password must contain at least one number")
	case !lowercaseRegexp.MatchString(password):
		return errors.New("Password must contain at least one lowercase letter")
	case !uppercaseRegexp.MatchString(password):
		return errors.New("Password must contain at least one uppercase letter")
	}
	return nil
}

// validateName проверяет валидность имени пользователя
func validateName(name string) error {
	if len(name) < 2 || len(name) > 24 {
		return errors.New("Name must be between 2 and 24 characters long")
	}

	if !nameRegexp.MatchString(name) {
		return errors.New("Name can only contain letters, spaces, and '-'")
	}

	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
