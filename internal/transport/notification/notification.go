package notification

import (
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/notification"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type NotificationService struct {
	uc notification.INotificationUsecase
}

func NewNotificationService(uc notification.INotificationUsecase) *NotificationService {
	return &NotificationService{uc: uc}
}

func (h *NotificationService) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	const op = "NotificationService.GetUserNotifications"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	offsetStr := vars["offset"]
	offset := 0
	var err error
    if offsetStr != "" {
        offset, err = strconv.Atoi(offsetStr)
        if err != nil {
            logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
            response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
            return
        }
    }
	
	notifications, err := h.uc.GetAllByUser(r.Context(), offset)
	if err != nil {
		logger.WithError(err).Error("failed to get notifications")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, notifications)
}

func (h *NotificationService) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	const op = "NotificationService.GetUnreadCount"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	count, err := h.uc.GetUnreadCount(r.Context())
	if err != nil {
		logger.WithError(err).Error("failed to get unread count")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, map[string]int{"unread_count": count})
}

func (h *NotificationService) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	const op = "NotificationService.MarkAsRead"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		logger.WithError(err).Error("invalid notification id")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	if err := h.uc.MarkAsRead(r.Context(), id); err != nil {
		logger.WithError(err).Error("failed to mark as read")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}