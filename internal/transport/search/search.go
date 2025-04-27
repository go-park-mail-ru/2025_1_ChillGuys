package search

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/suggestions"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/gorilla/mux"
)

//go:generate mockgen -source=search.go -destination=../../usecase/mocks/search_usecase_mock.go -package=mocks ISearchUsecase
type ISearchUsecase interface {
	SearchProductsByName(context.Context, dto.ProductNameResponse, int) ([]*models.Product, error)
	SearchCategoryByName(context.Context, dto.CategoryNameResponse) ([]*models.Category, error)
}

type SearchService struct {
	u ISearchUsecase
	s suggestions.ISuggestionsUsecase
}

func NewSearchService(u ISearchUsecase, s suggestions.ISuggestionsUsecase) *SearchService {
	return &SearchService{
		u: u,
		s: s,
	}
}

func (h *SearchService) Search(w http.ResponseWriter, r *http.Request) {
	const op = "SearchService.Search"
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

	// Чтение строки запроса
	var req dto.SearchReq
	if err := request.ParseData(r, &req); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	// Получение предложений по продуктам
	productResponse, err := h.s.GetProductSuggestions(r.Context(), req.SubString)
	if err != nil {
		logger.WithError(err).Error("failed to get product suggestions")
		response.HandleDomainError(r.Context(), w, err, "get product suggestions")
		return
	}

	// Получение продуктов по найденным предложениям
	products, err := h.u.SearchProductsByName(r.Context(), productResponse, offset)
	if err != nil {
		logger.WithError(err).Error("failed to search products by names")
		response.HandleDomainError(r.Context(), w, err, "search products by names")
		return
	}

	// Получение предложений по категориям
	categoryResponse, err := h.s.GetCategorySuggestions(r.Context(), req.SubString)
	if err != nil {
		logger.WithError(err).Error("failed to get category suggestions")
		response.HandleDomainError(r.Context(), w, err, "get category suggestions")
		return
	}

	// Получение категорий по найденным предложениям
	categories, err := h.u.SearchCategoryByName(r.Context(), categoryResponse)
	if err != nil {
		logger.WithError(err).Error("failed to search categories by names")
		response.HandleDomainError(r.Context(), w, err, "search categories by names")
		return
	}

	// Преобразуем продукты и категории
	convertToProductsResponse := dto.ConvertToProductsResponse(products)
	convertToCategoriesResponse := dto.ConvertToCategoriesResponse(categories)

	// Объединяем результаты в общий ответ
	searchResponse := dto.SearchResponse{
		Categories: convertToCategoriesResponse,
		Products:   convertToProductsResponse,
	}

	// Отправляем объединенный ответ
	response.SendJSONResponse(r.Context(), w, http.StatusOK, searchResponse)
}
