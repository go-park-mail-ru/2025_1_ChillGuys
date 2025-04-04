package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/cookie"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

var (
	emailRegexp     = regexp.MustCompile(`^\w+(\.\w*)*@\w+(\.\w{2,})+$`)
	nameRegexp      = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
	digitRegexp     = regexp.MustCompile(`[0-9]`)
	lowercaseRegexp = regexp.MustCompile(`[a-z]`)
	uppercaseRegexp = regexp.MustCompile(`[A-Z]`)
)

//go:generate mockgen -source=user.go -destination=../../usecase/mocks/user_usecase_mock.go -package=mocks
type IAuthUsecase interface {
	Register(context.Context, models.UserRegisterRequestDTO) (string, error)
	Login(context.Context, models.UserLoginRequestDTO) (string, error)
	Logout(context.Context) error
	GetMe(context.Context) (*models.User, error)
	UploadAvatar(context.Context, minio.FileDataType) (string, error)
}

type AuthHandler struct {
	u            IAuthUsecase
	log          *logrus.Logger
	minioService minio.Client
}

func NewAuthHandler(
	u IAuthUsecase,
	log *logrus.Logger,
	mS minio.Client,
) *AuthHandler {
	return &AuthHandler{
		u:            u,
		log:          log,
		minioService: mS,
	}
}

// @Summary		Login user
// @Description	Авторизация пользователя
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		models.UserLoginRequestDTO	true	"User credentials"
// @success		200		{}			-							"No Content"
// @Header			200		{string}	Set-Cookie					"Устанавливает JWT-токен в куки"
// @Failure		400		{object}	response.ErrorResponse		"Ошибка валидации"
// @Failure		401		{object}	response.ErrorResponse		"Неверные email или пароль"
// @Failure		500		{object}	response.ErrorResponse		"Внутренняя ошибка сервера"
// @Router			/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request models.UserLoginRequestDTO
	if errStatusCode, err := response.ParseData(r.Body, &request); err != nil {
		response.SendErrorResponse(w, errStatusCode, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	sanitizeUserLoginRequest(&request)

	if err := ValidateLoginCreds(request); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.u.Login(r.Context(), request)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	cookie.Cookie(w, token, string(cookie.Token))
	response.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary		Register user
// @Description	Создает нового пользователя, хеширует пароль и устанавливает JWT-токен в куки
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			input	body		models.UserRegisterRequestDTO	true	"Данные для регистрации"
// @success		200		{}			-								"No Content"
// @Header			200		{string}	Set-Cookie						"Устанавливает JWT-токен в куки"
// @Failure		400		{object}	response.ErrorResponse			"Некорректный запрос"
// @Failure		409		{object}	response.ErrorResponse			"Пользователь уже существует"
// @Failure		500		{object}	response.ErrorResponse			"Внутренняя ошибка сервера"
// @Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var request models.UserRegisterRequestDTO
	if errStatusCode, err := response.ParseData(r.Body, &request); err != nil {
		response.SendErrorResponse(w, errStatusCode, fmt.Sprintf("Failed to parse request body: %v", err))
		return
	}

	sanitizeUserRegistrationRequest(&request)

	if err := ValidateRegistrationCreds(request); err != nil {
		response.SendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.u.Register(r.Context(), request)
	if err != nil {
		response.HandleError(w, err)
		return
	}

	cookie.Cookie(w, token, string(cookie.Token))
	response.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary		Logout user
// @Description	Выход пользователя
// @Tags			auth
// @Security		TokenAuth
// @Success		200	{}			"No Content"
// @Failure		401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.u.Logout(r.Context()); err != nil {
		response.HandleError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     string(cookie.Token),
		Value:    "",
		Path:     "/",
		Expires:  time.Now().UTC().AddDate(0, 0, -1),
		HttpOnly: true,
		Secure:   true,
	})

	response.SendSuccessResponse(w, http.StatusOK, nil)
}

// @Summary		Get user info
// @Description	Получение информации о текущем пользователе
// @Tags			users
// @Security		TokenAuth
// @Produce		json
// @Success		200	{object}	models.User				"Информация о пользователе"
// @Failure		400	{object}	response.ErrorResponse	"Некорректный запрос"
// @Failure		401	{object}	response.ErrorResponse	"Пользователь не найден"
// @Failure		500	{object}	response.ErrorResponse	"Ошибка сервера"
// @Router			/users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, err := h.u.GetMe(r.Context())
	if err != nil {
		response.HandleError(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, user)
}

// @Summary		Upload avatar
// @Description	Загружает аватар пользователя
// @Tags			users
// @Accept			multipart/form-data
// @Produce		json
// @Param			file	formData	file					true	"Файл изображения"
// @Success		200		{object}	map[string]string		"URL загруженного аватара"
// @Failure		400		{object}	response.ErrorResponse	"Ошибка при обработке формы"
// @Failure		500		{object}	response.ErrorResponse	"Ошибка загрузки файла"
// @Security		TokenAuth
// @Router			/users/avatar [post]
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.log.Warnf("error parsing multipart form: %v", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "failed to parse form data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Warnf("error getting file from form: %v", err)
		response.SendErrorResponse(w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.log.Errorf("error reading file: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, "failed to read file")
		return
	}

	fileData := minio.FileDataType{
		FileName: header.Filename,
		Data:     fileBytes,
	}

	avatarURL, err := h.u.UploadAvatar(r.Context(), fileData)
	if err != nil {
		h.log.Errorf("upload error: %v", err)
		response.SendErrorResponse(w, http.StatusInternalServerError, "upload failed")
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, map[string]string{
		"imageURL": avatarURL,
	})
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

	if req.Surname.Valid {
		if err := validateName(req.Surname.String); err != nil {
			return err
		}
	}

	return nil
}

// validateEmail Функция валидации почты
func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email")
	}
	return nil
}

// validatePassword проверяет валидность пароля
func validatePassword(password string) error {
	switch {
	case len(password) < 8:
		return errors.New("password must be at least 8 characters")
	case !digitRegexp.MatchString(password):
		return errors.New("password must contain at least one number")
	case !lowercaseRegexp.MatchString(password):
		return errors.New("password must contain at least one lowercase letter")
	case !uppercaseRegexp.MatchString(password):
		return errors.New("password must contain at least one uppercase letter")
	}
	return nil
}

// validateName проверяет валидность имени пользователя
func validateName(name string) error {
	if len(name) < 2 || len(name) > 24 {
		return errors.New("name must be between 2 and 24 characters long")
	}

	if !nameRegexp.MatchString(name) {
		return errors.New("name can only contain letters, spaces, and '-'")
	}

	return nil
}

// sanitizeUserRegistrationRequest удаляет лишние пробелы из полей запроса регистрации пользователя
func sanitizeUserRegistrationRequest(req *models.UserRegisterRequestDTO) {
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	if req.Surname.Valid {
		req.Surname.String = strings.TrimSpace(req.Surname.String)
		req.Surname.Valid = req.Surname.String != ""
	}
}

// sanitizeUserLoginRequest удаляет лишние пробелы из полей запроса для логина пользователя
func sanitizeUserLoginRequest(req *models.UserLoginRequestDTO) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}
