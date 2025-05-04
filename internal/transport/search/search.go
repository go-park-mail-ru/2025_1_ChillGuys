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
	SearchProductsByName(context.Context, dto.ProductNameResponse) ([]*models.Product, error)
	SearchCategoryByName(context.Context, dto.CategoryNameResponse) ([]*models.Category, error)
	SearchProductsByNameWithFilterAndSort(
		ctx context.Context,
		req dto.ProductNameResponse,
		offset int,
		minPrice, maxPrice float64,
		minRating float32,
		sortOption models.SortOption,
	) ([]*models.Product, error)
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
	productResponse, err := h.s.GetProductSuggestionsOffset(r.Context(), req.CategoryID, req.SubString, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get product suggestions")
		response.HandleDomainError(r.Context(), w, err, "get product suggestions")
		return
	}

	// Получение продуктов по найденным предложениям
	products, err := h.u.SearchProductsByName(r.Context(), productResponse)
	if err != nil {
		logger.WithError(err).Error("failed to search products by names")
		response.HandleDomainError(r.Context(), w, err, "search products by names")
		return
	}

	var categories []*models.Category
	if !req.CategoryID.Valid {
		categoryResponse, err := h.s.GetCategorySuggestions(r.Context(), req.SubString)
		if err != nil {
			logger.WithError(err).Error("failed to get category suggestions")
			response.HandleDomainError(r.Context(), w, err, "get category suggestions")
			return
		}

		categories, err = h.u.SearchCategoryByName(r.Context(), categoryResponse)
		if err != nil {
			logger.WithError(err).Error("failed to search categories by names")
			response.HandleDomainError(r.Context(), w, err, "search categories by names")
			return
		}
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

func (h *SearchService) SearchWithFilterAndSort(w http.ResponseWriter, r *http.Request) {
	const op = "SearchService.SearchWithFilterAndSort"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	// Парсинг offset
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

	// Парсинг фильтров
	minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)
	minRating, _ := strconv.ParseFloat(r.URL.Query().Get("min_rating"), 32)

	// Парсинг параметра сортировки
	sortOption := models.SortOption(r.URL.Query().Get("sort"))
	switch sortOption {
	case models.SortByPriceAsc, models.SortByPriceDesc, models.SortByRatingAsc, models.SortByRatingDesc, models.SortByDefault:
		// допустимые значения
	default:
		sortOption = models.SortByDefault
	}

	// Чтение строки запроса
	var req dto.SearchReq
	if err := request.ParseData(r, &req); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	// Получение предложений по продуктам
	productResponse, err := h.s.GetProductSuggestionsOffset(r.Context(), req.CategoryID, req.SubString, offset)
	if err != nil {
		logger.WithError(err).Error("failed to get product suggestions")
		response.HandleDomainError(r.Context(), w, err, "get product suggestions")
		return
	}

	var categories []*models.Category
	if !req.CategoryID.Valid {
		categoryResponse, err := h.s.GetCategorySuggestions(r.Context(), req.SubString)
		if err != nil {
			logger.WithError(err).Error("failed to get category suggestions")
			response.HandleDomainError(r.Context(), w, err, "get category suggestions")
			return
		}

		categories, err = h.u.SearchCategoryByName(r.Context(), categoryResponse)
		if err != nil {
			logger.WithError(err).Error("failed to search categories by names")
			response.HandleDomainError(r.Context(), w, err, "search categories by names")
			return
		}
	}

	// Получение продуктов с фильтрацией и сортировкой
	products, err := h.u.SearchProductsByNameWithFilterAndSort(
		r.Context(),
		productResponse,
		offset,
		minPrice,
		maxPrice,
		float32(minRating),
		sortOption,
	)
	if err != nil {
		logger.WithError(err).Error("failed to search products by names")
		response.HandleDomainError(r.Context(), w, err, "search products by names")
		return
	}

	// Формирование ответа
	searchResponse := dto.SearchResponse{
		Categories: dto.ConvertToCategoriesResponse(categories),
		Products:   dto.ConvertToProductsResponse(products),
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, searchResponse)
}
