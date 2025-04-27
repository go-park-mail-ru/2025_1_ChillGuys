package review

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	gen "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/generated/review"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
)

type ReviewHandler struct {
	reviewClient gen.ReviewServiceClient
}

func NewReviewHandler(rc gen.ReviewServiceClient) *ReviewHandler{
	return &ReviewHandler{
		reviewClient: rc,
	}
}

func (h *ReviewHandler) Add(w http.ResponseWriter, r *http.Request){
	const op = "ReviewHandler.Add"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var addReq dto.AddReviewRequest
	if err := request.ParseData(r, &addReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.reviewClient.AddReview(r.Context(), dto.ConvertAddReviewRequestToGRPC(addReq))
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		logger.Error("add review")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, nil)
}

func (h *ReviewHandler) Get(w http.ResponseWriter, r *http.Request){
	const op = "ReviewHandler.Get"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var getReq dto.GetReviewRequest
	if err := request.ParseData(r, &getReq); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	reviews, err := h.reviewClient.GetReviews(r.Context(), dto.ConvertGetReviewRequestToGRPC(getReq))
	if err != nil {
		response.HandleGRPCError(r.Context(), w, err, op)
		logger.Error("get review")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, dto.ConvertGRPCToReviewsResponse(reviews))
}