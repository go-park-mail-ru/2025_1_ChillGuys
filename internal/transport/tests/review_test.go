package tests

import (
	"bytes"
	"encoding/json"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	genmock "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review/mocks"
	review "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/review/http"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReviewHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := genmock.NewMockReviewServiceClient(ctrl)

	handler := review.NewReviewHandler(mockClient)

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/reviews/get", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		handler.Get(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("get reviews failure", func(t *testing.T) {
		reqData := map[string]string{
			"product_id": "42",
		}
		body, _ := json.Marshal(reqData)

		// Mock the GetReviews method failure
		mockClient.EXPECT().GetReviews(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/reviews/get", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Get(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("successful get reviews", func(t *testing.T) {
		reqData := map[string]string{
			"product_id": "42",
		}
		body, _ := json.Marshal(reqData)

		// Mock successful GetReviews response
		mockClient.EXPECT().GetReviews(gomock.Any(), gomock.Any()).Return(&gen.GetReviewsResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/reviews/get", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Get(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}

func TestReviewHandler_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	t.Run("invalid body", func(t *testing.T) {
		mockClient := genmock.NewMockReviewServiceClient(ctrl)
		handler := review.NewReviewHandler(mockClient)

		req := httptest.NewRequest(http.MethodPost, "/reviews/add", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()

		handler.Add(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("add review failure", func(t *testing.T) {
		mockClient := genmock.NewMockReviewServiceClient(ctrl)
		handler := review.NewReviewHandler(mockClient)

		reqData := map[string]interface{}{
			"product_id": "42",
			"rating":     5,
			"text":       "Great product!",
		}
		body, _ := json.Marshal(reqData)

		// Mock the AddReview method failure
		mockClient.EXPECT().AddReview(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodPost, "/reviews/add", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Add(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("successful add review", func(t *testing.T) {
		mockClient := genmock.NewMockReviewServiceClient(ctrl)
		handler := review.NewReviewHandler(mockClient)

		reqData := map[string]interface{}{
			"product_id": "42",
			"rating":     5,
			"text":       "Great product!",
		}
		body, _ := json.Marshal(reqData)

		// Mock successful AddReview response
		mockClient.EXPECT().AddReview(gomock.Any(), gomock.Any()).Return(&gen.EmptyResponse{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/reviews/add", bytes.NewReader(body))
		w := httptest.NewRecorder()

		handler.Add(w, req)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})

	//t.Run("invalid rating value", func(t *testing.T) {
	//	mockClient := genmock.NewMockReviewServiceClient(ctrl)
	//	handler := review.NewReviewHandler(mockClient)
	//
	//	reqData := map[string]interface{}{
	//		"product_id": "42",
	//		"rating":     6, // invalid, should be 1-5
	//		"text":       "Great product!",
	//	}
	//	body, _ := json.Marshal(reqData)
	//
	//	req := httptest.NewRequest(http.MethodPost, "/reviews/add", bytes.NewReader(body))
	//	w := httptest.NewRecorder()
	//
	//	handler.Add(w, req)
	//
	//	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	//})
	//
	//t.Run("missing required fields", func(t *testing.T) {
	//	mockClient := genmock.NewMockReviewServiceClient(ctrl)
	//	handler := review.NewReviewHandler(mockClient)
	//
	//	reqData := map[string]interface{}{
	//		"product_id": "42",
	//		// missing rating
	//		"text": "Great product!",
	//	}
	//	body, _ := json.Marshal(reqData)
	//
	//	req := httptest.NewRequest(http.MethodPost, "/reviews/add", bytes.NewReader(body))
	//	w := httptest.NewRecorder()
	//
	//	handler.Add(w, req)
	//
	//	assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	//})
}
