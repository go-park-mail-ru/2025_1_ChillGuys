package tests

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/notification"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService_GetUserNotifications(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockINotificationUsecase(ctrl)
	srv := notification.NewNotificationService(mockUC)

	// Prepare context with logger
	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	// Helper to create request with mux vars
	makeRequest := func(offset string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"offset": offset})
		return req
	}

	t.Run("success with offset", func(t *testing.T) {
		offset := 5
		resp := dto.NotificationsListResponse{
			Total: 10,
		}

		mockUC.EXPECT().GetAllByUser(gomock.Any(), offset).Return(resp, nil)

		req := makeRequest(strconv.Itoa(offset))
		rr := httptest.NewRecorder()

		srv.GetUserNotifications(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"total":10`)
	})

	t.Run("success with empty offset (default 0)", func(t *testing.T) {
		offset := 0
		resp := dto.NotificationsListResponse{
			Total: 3,
		}
		mockUC.EXPECT().GetAllByUser(gomock.Any(), offset).Return(resp, nil)

		req := makeRequest("")
		rr := httptest.NewRecorder()

		srv.GetUserNotifications(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"total":3`)
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := makeRequest("notanint")
		rr := httptest.NewRecorder()

		srv.GetUserNotifications(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		offset := 0
		mockUC.EXPECT().GetAllByUser(gomock.Any(), offset).Return(dto.NotificationsListResponse{}, errors.New("fail"))

		req := makeRequest("")
		rr := httptest.NewRecorder()

		srv.GetUserNotifications(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockINotificationUsecase(ctrl)
	srv := notification.NewNotificationService(mockUC)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	t.Run("success", func(t *testing.T) {
		mockUC.EXPECT().GetUnreadCount(gomock.Any()).Return(7, nil)

		srv.GetUnreadCount(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"unread_count":7}`, rr.Body.String())
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUC.EXPECT().GetUnreadCount(gomock.Any()).Return(0, errors.New("fail"))

		rr := httptest.NewRecorder()
		srv.GetUnreadCount(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockINotificationUsecase(ctrl)
	srv := notification.NewNotificationService(mockUC)

	ctx := logctx.WithLogger(context.Background(), logrus.NewEntry(logrus.New()))

	makeRequest := func(id string) *http.Request {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": id})
		return req
	}

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		mockUC.EXPECT().MarkAsRead(gomock.Any(), id).Return(nil)

		req := makeRequest(id.String())
		rr := httptest.NewRecorder()

		srv.MarkAsRead(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "", rr.Body.String()) // json null response
	})

	t.Run("invalid UUID", func(t *testing.T) {
		req := makeRequest("not-a-uuid")
		rr := httptest.NewRecorder()

		srv.MarkAsRead(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		id := uuid.New()
		mockUC.EXPECT().MarkAsRead(gomock.Any(), id).Return(errors.New("fail"))

		req := makeRequest(id.String())
		rr := httptest.NewRecorder()

		srv.MarkAsRead(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
