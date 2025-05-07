package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/admin"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupTestAdmin(t *testing.T) (*mocks.MockIAdminUsecase, *admin.AdminService) {
	ctrl := gomock.NewController(t)
	mockUsecase := mocks.NewMockIAdminUsecase(ctrl)
	service := admin.NewAdminService(mockUsecase)
	return mockUsecase, service
}

func TestAdminService_GetPendingProducts(t *testing.T) {
	mockUsecase, service := setupTestAdmin(t)

	t.Run("success", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetPendingProducts(gomock.Any(), 0).
			Return(dto.ProductsResponse{Total: 1}, nil)

		req := httptest.NewRequest("GET", "/api/v1/admin/products/0", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "0"})
		w := httptest.NewRecorder()

		service.GetPendingProducts(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid offset", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/admin/products/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "invalid"})
		w := httptest.NewRecorder()

		service.GetPendingProducts(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetPendingProducts(gomock.Any(), 0).
			Return(dto.ProductsResponse{}, errors.New("fail"))

		req := httptest.NewRequest("GET", "/api/v1/admin/products/0", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "0"})
		w := httptest.NewRecorder()

		service.GetPendingProducts(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAdminService_UpdateProductStatus(t *testing.T) {
	mockUsecase, service := setupTestAdmin(t)

	t.Run("success", func(t *testing.T) {
		productID := uuid.New()
		reqBody := dto.UpdateProductStatusRequest{ProductID: productID, Update: 1}
		body, _ := json.Marshal(reqBody)

		mockUsecase.EXPECT().
			UpdateProductStatus(gomock.Any(), reqBody).
			Return(nil)

		req := httptest.NewRequest("POST", "/api/v1/admin/products/update-status", bytes.NewReader(body))
		w := httptest.NewRecorder()

		service.UpdateProductStatus(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/products/update-status", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		service.UpdateProductStatus(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		productID := uuid.New()
		reqBody := dto.UpdateProductStatusRequest{ProductID: productID, Update: 1}
		body, _ := json.Marshal(reqBody)

		mockUsecase.EXPECT().
			UpdateProductStatus(gomock.Any(), reqBody).
			Return(errors.New("fail"))

		req := httptest.NewRequest("POST", "/api/v1/admin/products/update-status", bytes.NewReader(body))
		w := httptest.NewRecorder()

		service.UpdateProductStatus(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAdminService_UpdateUserRole(t *testing.T) {
	mockUsecase, service := setupTestAdmin(t)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		reqBody := dto.UpdateUserRoleRequest{UserID: userID, Update: 1}
		body, _ := json.Marshal(reqBody)

		mockUsecase.EXPECT().
			UpdateUserRole(gomock.Any(), reqBody).
			Return(nil)

		req := httptest.NewRequest("POST", "/api/v1/admin/users/update-role", bytes.NewReader(body))
		w := httptest.NewRecorder()

		service.UpdateUserRole(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/users/update-role", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		service.UpdateUserRole(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("internal error", func(t *testing.T) {
		userID := uuid.New()
		reqBody := dto.UpdateUserRoleRequest{UserID: userID, Update: 1}
		body, _ := json.Marshal(reqBody)

		mockUsecase.EXPECT().
			UpdateUserRole(gomock.Any(), reqBody).
			Return(errs.ErrInternal)

		req := httptest.NewRequest("POST", "/api/v1/admin/users/update-role", bytes.NewReader(body))
		w := httptest.NewRecorder()

		service.UpdateUserRole(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAdminService_GetPendingUsers(t *testing.T) {
	mockUsecase, service := setupTestAdmin(t)

	t.Run("успешный запрос", func(t *testing.T) {
		expected := dto.UsersResponse{Total: 2}

		mockUsecase.EXPECT().
			GetPendingUsers(gomock.Any(), 0).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/0", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "0"})
		w := httptest.NewRecorder()

		service.GetPendingUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var actual dto.UsersResponse
		err := json.NewDecoder(resp.Body).Decode(&actual)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("ошибка парсинга offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/notanumber", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "notanumber"})
		w := httptest.NewRecorder()

		service.GetPendingUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("внутренняя ошибка usecase", func(t *testing.T) {
		mockUsecase.EXPECT().
			GetPendingUsers(gomock.Any(), 5).
			Return(dto.UsersResponse{}, errors.New("db fail"))

		req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/5", nil)
		req = mux.SetURLVars(req, map[string]string{"offset": "5"})
		w := httptest.NewRecorder()

		service.GetPendingUsers(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
