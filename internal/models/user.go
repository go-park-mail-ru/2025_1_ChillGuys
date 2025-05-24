package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

const (
	RoleAdmin   UserRole = "admin"   // админ
	RoleBuyer   UserRole = "buyer"   // покупатель
	RoleSeller  UserRole = "seller"  // продавец
	RolePending UserRole = "pending" // в ожидании
	RoleWarehouseman UserRole = "warehouseman"
)

// UserRole представляет роль пользователя в системе
type UserRole string

// String возвращает строковое представление роли
func (r UserRole) String() string {
	return string(r)
}

// ParseUserRole преобразует строку в UserRole
func ParseUserRole(role string) (UserRole, error) {
	switch role {
	case "admin":
		return RoleAdmin, nil
	case "buyer":
		return RoleBuyer, nil
	case "seller":
		return RoleSeller, nil
	case "pending":
		return RolePending, nil
	case "warehouseman":
		return RoleWarehouseman, nil
	default:
		return RolePending, fmt.Errorf("unknown user role: %s", role)
	}
}

// Scan реализует интерфейс sql.Scanner для чтения из БД
func (r *UserRole) Scan(value interface{}) error {
	if value == nil {
		*r = RolePending
		return nil
	}

	var roleStr string

	switch v := value.(type) {
	case string:
		roleStr = v
	case []byte:
		roleStr = string(v)
	default:
		return fmt.Errorf("failed to scan UserRole: unsupported type %T", value)
	}

	role, err := ParseUserRole(roleStr)
	if err != nil {
		return err
	}
	*r = role
	return nil
}

// Value реализует интерфейс driver.Valuer для записи в БД
func (r UserRole) Value() (driver.Value, error) {
	return r.String(), nil
}

type User struct {
	ID          uuid.UUID   `json:"id"`
	Email       string      `json:"email"`
	Name        string      `json:"name"`
	Surname     null.String `json:"surname" swaggertype:"primitive,string"`
	ImageURL    null.String `json:"imageURL" swaggertype:"primitive,string"`
	PhoneNumber null.String `json:"phoneNumber,omitempty" swaggertype:"primitive,string"`
	Role        UserRole    `json:"role"`
	Seller      *Seller     `json:"seller,omitempty"`
}

type Seller struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
}

type UpdateUserDB struct {
	Name        string
	Surname     null.String
	PhoneNumber null.String
}

type UserDB struct {
	ID           uuid.UUID
	Email        string
	Name         string
	Surname      null.String
	ImageURL     null.String
	PhoneNumber  null.String
	PasswordHash []byte
	Role         UserRole
}

func (u *UserDB) ConvertToUser() *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Name:        u.Name,
		Surname:     u.Surname,
		ImageURL:    u.ImageURL,
		PhoneNumber: u.PhoneNumber,
		Role:		 u.Role,
	}
}
