package dto

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/google/uuid"
)

type UpdateUserRoleRequest struct{
	UserID uuid.UUID 	  `json:"userID"`
	Update int 			  `json:"update"`
}

type UpdateProductStatusRequest struct{
	ProductID uuid.UUID   `json:"productID"`
	Update int 			  `json:"update"`
}

type BriefUser struct {
    ID         uuid.UUID `json:"id"`
    Email      string    `json:"email"`
    Name       string    `json:"name"`
    Surname    string    `json:"surname,omitempty"`
    ImageURL   string    `json:"image,omitempty"`
    Phone      string    `json:"phone,omitempty"`
    Role       string    `json:"role"`
    SellerInfo *SellerInfo `json:"seller_info,omitempty"`
}

type SellerInfo struct {
    Title       string `json:"title"`
    Description string `json:"description,omitempty"`
}

func ConvertToBriefUser(user *models.User) BriefUser {
    briefUser := BriefUser{
        ID:       user.ID,
        Email:    user.Email,
        Name:     user.Name,
        Surname:  user.Surname.String,
        ImageURL: user.ImageURL.String,
        Phone:    user.PhoneNumber.String,
        Role:     user.Role.String(),
    }

    if user.Seller != nil {
        briefUser.SellerInfo = &SellerInfo{
            Title:       user.Seller.Title,
            Description: user.Seller.Description,
        }
    }

    return briefUser
}

type UsersResponse struct {
    Total int         `json:"total"`
    Users []BriefUser `json:"users"`
}

func ConvertToUsersResponse(users []*models.User) UsersResponse {
    briefUsers := make([]BriefUser, 0, len(users))
    for _, user := range users {
        briefUsers = append(briefUsers, ConvertToBriefUser(user))
    }

    return UsersResponse{
        Total: len(briefUsers),
        Users: briefUsers,
    }
}