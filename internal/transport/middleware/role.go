package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/domains"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/sirupsen/logrus"
)

// RoleMiddleware создает middleware для проверки роли
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := logctx.GetLogger(ctx)

			// Получаем роль из контекста
			role, ok := ctx.Value(domains.RoleKey{}).(string)
			if !ok {
				logger.Error("Role not found in context")
				response.SendJSONError(r.Context(), w, http.StatusInternalServerError, "Role not found in context")
				return
			}

			for _, allowedRole := range allowedRoles {
                if role == allowedRole {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            logger.WithFields(logrus.Fields{
                "required_roles": allowedRoles,
                "user_role":      role,
                "path":           r.URL.Path,
            }).Warn("Access denied")
            
            response.SendJSONError(ctx, w, http.StatusForbidden, "Insufficient permissions")
		})
	}
}