package suggestions

import (
	"context"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/guregu/null"
	"github.com/mailru/easyjson"
	"net/http"
)

//go:generate mockgen -source=suggestions.go -destination=../../usecase/mocks/suggestions_usecase_mock.go -package=mocks ISuggestionsUsecase
type ISuggestionsUsecase interface {
	GetCategorySuggestions(context.Context, string) (dto.CategoryNameResponse, error)
	GetProductSuggestions(context.Context, null.String, string) (dto.ProductNameResponse, error)
}

type SuggestionsService struct {
	u ISuggestionsUsecase
}

func NewSuggestionsService(u ISuggestionsUsecase) *SuggestionsService {
	return &SuggestionsService{
		u: u,
	}
}

func (h *SuggestionsService) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	const op = "SuggestionsService.GetSuggestions"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.SuggestionsReq
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	var categoryResponse dto.CategoryNameResponse
	var productResponse dto.ProductNameResponse
	var err error

	// Если категория не указана - получаем и категории и продукты
	if !req.CategoryID.Valid {
		categoryResponse, err = h.u.GetCategorySuggestions(r.Context(), req.SubString)
		if err != nil {
			logger.WithError(err).Error("failed to get category suggestions")
			response.HandleDomainError(r.Context(), w, err, "get category suggestions")
			return
		}

		productResponse, err = h.u.GetProductSuggestions(r.Context(), req.CategoryID, req.SubString)
		if err != nil {
			logger.WithError(err).Error("failed to get product suggestions")
			response.HandleDomainError(r.Context(), w, err, "get product suggestions")
			return
		}
	} else {
		// Если категория указана - получаем только продукты для этой категории
		productResponse, err = h.u.GetProductSuggestions(r.Context(), req.CategoryID, req.SubString)
		if err != nil {
			logger.WithError(err).Error("failed to get product suggestions")
			response.HandleDomainError(r.Context(), w, err, "get product suggestions")
			return
		}
		// Возвращаем пустой список категорий
		categoryResponse = dto.CategoryNameResponse{CategoriesNames: []models.CategorySuggestion{}}
	}

	combined := dto.CombinedSuggestionsResponse{
		Categories: categoryResponse.CategoriesNames,
		Products:   productResponse.ProductNames,
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, combined)
}
