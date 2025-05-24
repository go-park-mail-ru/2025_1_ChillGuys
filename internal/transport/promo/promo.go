package promo

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/gorilla/mux"
)

//go:generate mockgen -source=promo.go -destination=../../usecase/mocks/promo_usecase_mock.go -package=mocks IPromoUsecase
type IPromoUsecase interface {
	CreatePromo(ctx context.Context, req dto.CreatePromoRequest) (dto.PromoResponse, error)
	GetAllPromos(ctx context.Context, offset int) (dto.PromosResponse, error)
	CheckPromoCode(ctx context.Context, code string) (dto.PromoValidityResponse, error)
}

type PromoService struct {
	uc IPromoUsecase
}

func NewPromoService(uc IPromoUsecase) *PromoService {
	return &PromoService{uc: uc}
}

func (h *PromoService) Create(w http.ResponseWriter, r *http.Request) {
	const op = "PromoService.CreatePromo"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	var req dto.CreatePromoRequest
	if err := request.ParseData(r, &req); err != nil {
		logger.WithError(err).Error("parse request data")
		response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
		return
	}

	resp, err := h.uc.CreatePromo(r.Context(), req)
	if err != nil {
		logger.WithError(err).Error("create promo")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, resp)
}

func (h *PromoService) GetAll(w http.ResponseWriter, r *http.Request) {
	const op = "PromoService.GetAllPromos"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	offsetStr := vars["offset"]
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
		response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
		return
	}

	promos, err := h.uc.GetAllPromos(r.Context(), offset)
	if err != nil {
		logger.WithError(err).Error("get all promos")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, promos)
}

func (h *PromoService) CheckPromoCode(w http.ResponseWriter, r *http.Request) {
    const op = "PromoService.GetAllPromos"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)
    
    var req dto.CheckPromoRequest
    if err := request.ParseData(r, &req); err != nil {
		logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
        return
    }
    
    result, err := h.uc.CheckPromoCode(r.Context(), req.Code)
    if err != nil {
		logger.WithError(err).Error("get promo")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }
    
    response.SendJSONResponse(r.Context(), w, http.StatusOK, result)
}