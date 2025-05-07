package tests

import (
	"bytes"
	"encoding/json"
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
}
