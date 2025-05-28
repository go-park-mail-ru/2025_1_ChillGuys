package tests

import (
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/recommendation"
	mockRecommendation "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRecommendations_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockRecommendation.NewMockIRecommendationUsecase(ctrl)
	handler := recommendation.NewRecommendationService(mockUsecase)

	productID := uuid.New()
	expectedProducts := []*models.Product{
		{
			ID:              uuid.New(),
			Name:            "Test Product",
			PreviewImageURL: "/img/test.png",
			Price:           100.0,
			PriceDiscount:   90.0,
			Quantity:        10,
			Rating:          4.5,
			ReviewsCount:    5,
		},
	}

	mockUsecase.EXPECT().
		GetRecommendations(gomock.Any(), productID).
		Return(expectedProducts, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recommendations/"+productID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

	rec := httptest.NewRecorder()
	handler.GetRecommendations(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var respBody map[string]interface{}
	err := json.NewDecoder(rec.Body).Decode(&respBody)
	require.NoError(t, err)

	assert.Equal(t, float64(1), respBody["total"])
	assert.NotNil(t, respBody["products"])
}

func TestGetRecommendations_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockRecommendation.NewMockIRecommendationUsecase(ctrl)
	handler := recommendation.NewRecommendationService(mockUsecase)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recommendations/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid-uuid"})

	rec := httptest.NewRecorder()
	handler.GetRecommendations(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetRecommendations_UsecaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mockRecommendation.NewMockIRecommendationUsecase(ctrl)
	handler := recommendation.NewRecommendationService(mockUsecase)

	productID := uuid.New()

	mockUsecase.EXPECT().
		GetRecommendations(gomock.Any(), productID).
		Return(nil, errors.New("something went wrong"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recommendations/"+productID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": productID.String()})

	rec := httptest.NewRecorder()
	handler.GetRecommendations(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
