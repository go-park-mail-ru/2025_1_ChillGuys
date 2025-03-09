package transport

import (
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/jwt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

//go:generate mockgen -source=user.go -destination=../repository/mocks/user_repo_mock.go -package=mocks IUserRepository

var (
	emailRegexp     = regexp.MustCompile(`^\w+(\.\w*)*@\w+(\.\w{2,})+$`)
	digitRegexp     = regexp.MustCompile(`[0-9]`)
	lowercaseRegexp = regexp.MustCompile(`[a-z]`)
	uppercaseRegexp = regexp.MustCompile(`[A-Z]`)
)

type IUserRepository interface {
	CreateUser(user models.UserRepo) error
	GetUserByEmail(email string) (*models.UserRepo, error)
	GetUserByID(id uuid.UUID) (*models.UserRepo, error)
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

// @Summary Login user
// @Description Авторизация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.UserLoginRequestDTO true "User credentials"
// @success 200 {} - "No Content"
// @Header 200 {string} Set-Cookie "Устанавливает JWT-токен в куки"
// @Failure 400 {object} utils.ErrorResponse "Ошибка валидации"
// @Failure 401 {object} utils.ErrorResponse "Неверные email или пароль"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request models.UserLoginRequestDTO
	if errStatusCode, errMessage := utils.ParseData(r.Body, &request); errStatusCode != 0 && errMessage != "" {
		utils.SendErrorResponse(w, errStatusCode, errMessage)
		return
	}

	if err := validateEmail(request.Email); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid email")
		return
	}

	if err := validatePassword(request.Password); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid password: %v", err))
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

// @Summary Register user
// @Description Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.UserRegisterRequestDTO true "Данные для регистрации"
// @success 200 {} - "No Content"
// @Header 200 {string} Set-Cookie "Устанавливает JWT-токен в куки"
// @Failure 400 {object} utils.ErrorResponse "Некорректный запрос"
// @Failure 409 {object} utils.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} utils.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request models.UserRegisterRequestDTO
	if errStatusCode, errMessage := utils.ParseData(r.Body, &request); errStatusCode != 0 && errMessage != "" {
		utils.SendErrorResponse(w, errStatusCode, errMessage)
		return
	}

	if err := validateEmail(request.Email); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid email")
		return
	}

	if err := validatePassword(request.Password); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid password: %v", err))
		return
	}

	if err := validateName(request.Name); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid name")
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

	userRepo := models.UserRepo{
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

// @Summary Logout user
// @Description Выход пользователя
// @Tags auth
// @Security ApiKeyAuth
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, isExist := r.Context().Value("userID").(string)
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

// @Summary Get user info
// @Description Получение информации о текущем пользователе
// @Tags users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} models.User "Информация о пользователе"
// @Failure 401 {object} utils.ErrorResponse "Неверный токен"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userIDStr, ok := r.Context().Value("userID").(string)
	if !ok {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "User id not found")
		return
	}

	version, ok := r.Context().Value("userVersion").(int)
	if !ok {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "User version not found")
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

	if !userRepo.IsVersionValid(version) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	user := userRepo.ConvertToUser()
	utils.SendSuccessResponse(w, http.StatusOK, user)
}

// validateEmail Функция валидации почты
func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePassword проверяет валидность пароля
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !digitRegexp.MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	if !lowercaseRegexp.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercaseRegexp.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	return nil
}

// validateName проверяет валидность имени пользователя
func validateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
