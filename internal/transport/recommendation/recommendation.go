package recommendation

import (
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/recommendation"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type RecommendationServise struct {
	r recommendation.IRecommendationUsecase
}

func NewRecommendationService(r recommendation.IRecommendationUsecase) *RecommendationServise {
	return &RecommendationServise{
		r: r,
	}
}

func (h *RecommendationServise) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.GetAllProducts"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	productID := vars["id"]

	parseProductID, err := uuid.Parse(productID)
	if err != nil {
		logger.WithError(err).WithField("productID", parseProductID).Error("parse productID")
		response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
		return
	}

	recommendations, err := h.r.GetRecommendations(r.Context(), parseProductID)
	if err != nil {
		logger.WithError(err).WithField("productID", parseProductID).Error("get recommendations")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	recommendationsResp := dto.ConvertToProductsResponse(recommendations)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, recommendationsResp)
}
