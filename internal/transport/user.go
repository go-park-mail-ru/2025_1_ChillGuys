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
)

//go:generate mockgen -source=auth.go -destination=../repository/mocks/user_repo_mock.go -package=mocks IUserRepository

type IUserRepository interface {
	CreateUser(user models.UserRepo) error
	GetUserByEmail(email string) (*models.UserRepo, error)
	IncrementUserVersion(userID string) error
}

type AuthHandler struct {
	repo IUserRepository
	log  *logrus.Logger
}

func NewAuthHandler(repo IUserRepository, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		repo: repo,
		log:  log,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	// Парсим
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

	// Получаем данные пользоваетеля из бд
	userRepo, err := h.repo.GetUserByEmail(request.Email)
	if err != nil {
		h.log.Warn(err.Error())
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	// Проверяем совпадение пароли
	if err := bcrypt.CompareHashAndPassword(userRepo.PasswordHash, []byte(request.Password)); err != nil {
		utils.SendErrorResponse(w, http.StatusUnauthorized, "Invalid password or email")
		return
	}

	token, err := jwt.CreateJWT(userRepo.ID.String(), userRepo.Version)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, &models.UserResponseDTO{
		Token: token,
	})
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

	if err := validateName(request.Surname); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid surname")
		return
	}

	// Создаём хэш пароля
	passwordHash, err := generatePasswordHash(request.Password)
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

	// Создаём jwt
	token, err := jwt.CreateJWT(userRepo.ID.String(), 1)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, &models.UserResponseDTO{
		Token: token,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	userID, isExist := r.Context().Value("userID").(string)
	if !isExist {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "user id not found")
		return
	}

	if err := h.repo.IncrementUserVersion(userID); err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, nil)
}

var (
	passwordRegexp = regexp.MustCompile(`^[a-zA-Z0-9]{8,}$`)
	emailRegexp    = regexp.MustCompile(`^[a-z0-9]+@[a-z0-9]+\.[a-z]{2,4}$`)
	nameRegexp     = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
)

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

// generatePasswordHash Генерация хэша пароля
func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
}
