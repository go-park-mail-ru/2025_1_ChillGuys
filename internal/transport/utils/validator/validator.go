package validator

import (
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"regexp"
	"strings"
)

var (
	emailRegexp       = regexp.MustCompile(`^\w+(\.\w*)*@\w+(\.\w{2,})+$`)
	nameRegexp        = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ\s-]+$`)
	digitRegexp       = regexp.MustCompile(`[0-9]`)
	lowercaseRegexp   = regexp.MustCompile(`[a-z]`)
	uppercaseRegexp   = regexp.MustCompile(`[A-Z]`)
	phoneNumberRegexp = regexp.MustCompile(`^(?:\+?\d{1,3})?[\s\-]?\(?\d{1,4}\)?[\s\-]?\d{1,4}[\s\-]?\d{1,4}$`)
)

func ValidateLoginCreds(req dto.UserLoginRequestDTO) error {
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// ValidateRegistrationCreds проверяет корректность данных при регистрации
func ValidateRegistrationCreds(req dto.UserRegisterRequestDTO) error {
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

func ValidateUserUpdateProfileCreds(req dto.UpdateUserProfileRequestDTO) error {

	if req.Email.Valid && strings.TrimSpace(req.Email.String) != "" {
		if err := validateEmail(req.Email.String); err != nil {
			return err
		}
	}

	if req.Password.Valid && strings.TrimSpace(req.Password.String) != "" {
		if err := validatePassword(req.Password.String); err != nil {
			return err
		}
	}

	if req.Name.Valid && strings.TrimSpace(req.Name.String) != "" {
		if err := validateName(req.Name.String); err != nil {
			return err
		}
	}

	if req.Surname.Valid && strings.TrimSpace(req.Surname.String) != "" {
		if err := validateName(req.Surname.String); err != nil {
			return err
		}
	}

	if req.PhoneNumber.Valid && strings.TrimSpace(req.PhoneNumber.String) != "" {
		if err := validatePhoneNumber(req.PhoneNumber.String); err != nil {
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

func validatePhoneNumber(phone string) error {
	if !phoneNumberRegexp.MatchString(phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}

// SanitizeUserRegistrationRequest удаляет лишние пробелы из полей запроса регистрации пользователя
func SanitizeUserRegistrationRequest(req *dto.UserRegisterRequestDTO) {
	req.Email = strings.TrimSpace(req.Email)
	req.Name = strings.TrimSpace(req.Name)
	req.Password = strings.TrimSpace(req.Password)
	if req.Surname.Valid {
		req.Surname.String = strings.TrimSpace(req.Surname.String)
		req.Surname.Valid = req.Surname.String != ""
	}
}

// SanitizeUserLoginRequest удаляет лишние пробелы из полей запроса для логина пользователя
func SanitizeUserLoginRequest(req *dto.UserLoginRequestDTO) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
}

func SanitizeUserProfileUpdateRequest(req *dto.UpdateUserProfileRequestDTO) {
	if req.Email.Valid {
		req.Email.String = strings.TrimSpace(req.Email.String)
	}

	if req.Name.Valid {
		req.Name.String = strings.TrimSpace(req.Name.String)
		req.Name.Valid = req.Name.String != ""
	}

	if req.Password.Valid {
		req.Password.String = strings.TrimSpace(req.Password.String)
	}

	if req.Surname.Valid {
		req.Surname.String = strings.TrimSpace(req.Surname.String)
		req.Surname.Valid = req.Surname.String != ""
	}

	if req.PhoneNumber.Valid {
		req.PhoneNumber.String = strings.TrimSpace(req.PhoneNumber.String)
		req.PhoneNumber.Valid = req.PhoneNumber.String != ""
	}
}
