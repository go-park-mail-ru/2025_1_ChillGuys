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
	"time"
)

//go:generate mockgen -source=user.go -destination=../repository/mocks/user_repo_mock.go -package=mocks IUserRepository

var (
	passwordRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	emailRegexp    = regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,4}$`)
	nameRegexp     = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Парсим запрос
	var request models.UserLoginRequestDTO
	if errStatusCode, errMessage := utils.ParseData(r.Body, &request); errStatusCode != 0 && errMessage != "" {
		utils.SendErrorResponse(w, errStatusCode, errMessage)
		return
	}

	// Валидация
	if err := validateEmail(request.Email); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid email")
		return
	}

	if err := validatePassword(request.Password); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid password: %v", err))
		return
	}

	// Получаем данные пользователя из базы данных
	userRepo, err := h.repo.GetUserByEmail(request.Email)
	if err != nil {
		h.log.Warn(err.Error())
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	// Проверяем совпадение пароля
	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(request.Password)); err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	// Генерация токена
	token, err := h.token.CreateJWT(userRepo.ID.String(), userRepo.Version)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Cookie(w, token, string(utils.Token))

	// Отправляем успешный ответ с токеном и версией
	utils.SendSuccessResponse(w, http.StatusOK, nil)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Парсим
	var request models.UserRegisterRequestDTO
	if errStatusCode, errMessage := utils.ParseData(r.Body, &request); errStatusCode != 0 && errMessage != "" {
		utils.SendErrorResponse(w, errStatusCode, errMessage)
		return
	}

	// Валидация
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

	// Создаём хэш пароля
	passwordHash, err := GeneratePasswordHash(request.Password)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Проверка существования пользователя в бд
	existedUser, _ := h.repo.GetUserByEmail(request.Email)
	if existedUser != nil {
		utils.SendErrorResponse(w, http.StatusConflict, "User already exists")
		return
	}

	// Запись в бд
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

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Получаем ID и версию из контекста (устанавливается в JWTMiddleware)
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

	//Проверяем, совпадает ли версия пользователя
	if !userRepo.IsVersionValid(version) {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Token is invalid or expired")
		return
	}

	// Преобразуем userRepo в user
	user := userRepo.ConvertToUser()

	// Отправляем успешный ответ с краткой информацией о пользователе
	utils.SendSuccessResponse(w, http.StatusOK, user)
}

// validateEmail Функция валидации почты
func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePassword Функция для проверки валидности пароля
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if !passwordRegexp.MatchString(password) {
		return errors.New("password must contain at least one letter and one number")
	}
	return nil
}

// validateName Функция валидации имени пользователя
func validateName(name string) error {
	if len(name) < 2 || len(name) > 50 {
		return errors.New("name must be between 2 and 50 characters long")
	}

	if !nameRegexp.MatchString(name) {
		return errors.New("name can only contain letters, spaces, and '-'")
	}

	return nil
}

// GeneratePasswordHash Генерация хэша пароля
func GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
