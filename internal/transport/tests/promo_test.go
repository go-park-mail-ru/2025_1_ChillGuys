package tests

import (
	"bytes"
	"github.com/google/uuid"

	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/promo"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
)

func stripMono(t time.Time) time.Time {
	return t.Round(0)
}

func TestPromoService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockIPromoUsecase(ctrl)
	service := promo.NewPromoService(mockUC)

	validID := uuid.New()

	validReq := dto.CreatePromoRequest{
		Code:      "PROMO1",
		Percent:   10,
		StartDate: stripMono(time.Now()),
		EndDate:   stripMono(time.Now().Add(24 * time.Hour)),
	}

	validResp := dto.PromoResponse{
		ID:        validID, // просто для примера, в реале uuid
		Code:      validReq.Code,
		Percent:   validReq.Percent,
		StartDate: validReq.StartDate,
		EndDate:   validReq.EndDate,
	}

	body, _ := easyjson.Marshal(&validReq)
	req := httptest.NewRequest(http.MethodPost, "/promo", bytes.NewReader(body))
	w := httptest.NewRecorder()

	mockUC.EXPECT().CreatePromo(gomock.Any(), validReq).Return(validResp, nil)

	service.Create(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	if !strings.Contains(w.Body.String(), validReq.Code) {
		t.Errorf("response body does not contain promo code")
	}

	// неверный json
	badReq := httptest.NewRequest(http.MethodPost, "/promo", bytes.NewReader([]byte("bad json")))
	w = httptest.NewRecorder()
	service.Create(w, badReq)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on bad json, got %d", http.StatusBadRequest, w.Code)
	}

	// ошибка usecase
	req = httptest.NewRequest(http.MethodPost, "/promo", bytes.NewReader(body))
	w = httptest.NewRecorder()

	mockUC.EXPECT().
		CreatePromo(gomock.Any(), gomock.AssignableToTypeOf(dto.CreatePromoRequest{})).
		Return(dto.PromoResponse{}, errors.New("fail"))

	service.Create(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on usecase error, got %d", http.StatusInternalServerError, w.Code)
	}

}

func TestPromoService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	validID, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")

	mockUC := mocks.NewMockIPromoUsecase(ctrl)
	service := promo.NewPromoService(mockUC)

	// валидный запрос с offset=0
	req := httptest.NewRequest(http.MethodGet, "/promo/0", nil)
	w := httptest.NewRecorder()

	// mux.Vars нужен для получения параметра offset
	req = mux.SetURLVars(req, map[string]string{"offset": "0"})

	expectedResp := dto.PromosResponse{
		Total: 1,
		Promos: []dto.PromoResponse{
			{
				ID:      validID,
				Code:    "PROMO1",
				Percent: 10,
			},
		},
	}

	mockUC.EXPECT().GetAllPromos(gomock.Any(), 0).Return(expectedResp, nil)

	service.GetAll(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "PROMO1") {
		t.Errorf("response body does not contain promo code")
	}

	// неверный offset
	req = httptest.NewRequest(http.MethodGet, "/promo/bad", nil)
	req = mux.SetURLVars(req, map[string]string{"offset": "bad"})
	w = httptest.NewRecorder()

	service.GetAll(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on bad offset, got %d", http.StatusInternalServerError, w.Code)
	}

	// ошибка usecase
	req = httptest.NewRequest(http.MethodGet, "/promo/0", nil)
	req = mux.SetURLVars(req, map[string]string{"offset": "0"})
	w = httptest.NewRecorder()
	mockUC.EXPECT().GetAllPromos(gomock.Any(), 0).Return(dto.PromosResponse{}, errors.New("fail"))
	service.GetAll(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on usecase error, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestPromoService_CheckPromoCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := mocks.NewMockIPromoUsecase(ctrl)
	service := promo.NewPromoService(mockUC)

	// валидный запрос
	validReq := dto.CheckPromoRequest{
		Code: "PROMO1",
	}

	body, _ := easyjson.Marshal(&validReq)
	req := httptest.NewRequest(http.MethodPost, "/promo/check", bytes.NewReader(body))
	w := httptest.NewRecorder()

	expectedResp := dto.PromoValidityResponse{
		IsValid: true,
		Percent: 10,
	}

	mockUC.EXPECT().CheckPromoCode(gomock.Any(), validReq.Code).Return(expectedResp, nil)

	service.CheckPromoCode(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "true") {
		t.Errorf("response body does not contain is_valid true")
	}

	// неверный json
	badReq := httptest.NewRequest(http.MethodPost, "/promo/check", bytes.NewReader([]byte("bad json")))
	w = httptest.NewRecorder()
	service.CheckPromoCode(w, badReq)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on bad json, got %d", http.StatusInternalServerError, w.Code)
	}

	// ошибка usecase
	req = httptest.NewRequest(http.MethodPost, "/promo/check", bytes.NewReader(body))
	w = httptest.NewRecorder()
	mockUC.EXPECT().CheckPromoCode(gomock.Any(), validReq.Code).Return(dto.PromoValidityResponse{}, errors.New("fail"))
	service.CheckPromoCode(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d on usecase error, got %d", http.StatusInternalServerError, w.Code)
	}
}
