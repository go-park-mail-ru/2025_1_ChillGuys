package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/guregu/null"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/search"
	usecasemocks "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestSearchWithFilterAndSort_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearchUC := usecasemocks.NewMockISearchUsecase(ctrl)
	mockSuggestUC := usecasemocks.NewMockISuggestionsUsecase(ctrl)

	handler := search.NewSearchService(mockSearchUC, mockSuggestUC)

	reqBody := dto.SearchReq{
		CategoryID: null.String{}, // пустой — будем использовать подсказки
		SubString:  "phone",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/search/0?sort=price_asc&min_price=100&max_price=1000&min_rating=3", bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"offset": "0"})
	logger := logrus.NewEntry(logrus.New()) // или какой логгер вы используете
	ctx := logctx.WithLogger(context.Background(), logger)
	req = req.WithContext(ctx)

	// Category suggestions
	categoryResp := dto.CategoryNameResponse{
		CategoriesNames: []models.CategorySuggestion{{Name: "Phones"}},
	}

	mockSuggestUC.EXPECT().GetCategorySuggestions(gomock.Any(), "phone").Return(categoryResp, nil)

	// Category search
	categories := []*models.Category{{Name: "Phones"}}
	mockSearchUC.EXPECT().SearchCategoryByName(gomock.Any(), categoryResp).Return(categories, nil)

	// Products search
	mockSearchUC.EXPECT().SearchProductsByNameWithFilterAndSort(
		gomock.Any(), null.String{}, "phone", 0, 100.0, 1000.0, float32(3.0), models.SortByPriceAsc,
	).Return([]*models.Product{}, nil)

	rr := httptest.NewRecorder()
	handler.SearchWithFilterAndSort(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"categories"`)
}

func TestSearchWithFilterAndSort_ParseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler := search.NewSearchService(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/search/abc", bytes.NewReader([]byte(`invalid json`)))
	req = mux.SetURLVars(req, map[string]string{"offset": "abc"})
	logger := logrus.NewEntry(logrus.New()) // или какой логгер вы используете
	ctx := logctx.WithLogger(context.Background(), logger)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.SearchWithFilterAndSort(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "failed to parse request body")
}

func TestSearchWithFilterAndSort_GetCategorySuggestionsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSearchUC := usecasemocks.NewMockISearchUsecase(ctrl)
	mockSuggestUC := usecasemocks.NewMockISuggestionsUsecase(ctrl)

	handler := search.NewSearchService(mockSearchUC, mockSuggestUC)

	reqBody := dto.SearchReq{
		CategoryID: null.String{},
		SubString:  "tablet",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/search/0", bytes.NewReader(jsonBody))
	req = mux.SetURLVars(req, map[string]string{"offset": "0"})
	logger := logrus.NewEntry(logrus.New()) // или какой логгер вы используете
	ctx := logctx.WithLogger(context.Background(), logger)
	req = req.WithContext(ctx)

	mockSuggestUC.EXPECT().GetCategorySuggestions(gomock.Any(), "tablet").Return(dto.CategoryNameResponse{}, errors.New("suggestion error"))

	rr := httptest.NewRecorder()
	handler.SearchWithFilterAndSort(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "suggestion error")
}
