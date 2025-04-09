package helpers

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
)

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
    userIDStr, isExist := ctx.Value(domains.UserIDKey).(string)
    if !isExist {
        return uuid.Nil, errs.NewNotFoundError("user not found")
    }
    
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        return uuid.Nil, errs.ErrInvalidID
    }
    
    return userID, nil
}