package admin

import (
    "context"
    "net/http"
    "strconv"

    "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
    "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
    "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
    "github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
    "github.com/gorilla/mux"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
)

//go:generate mockgen -source=admin.go -destination=../../usecase/mocks/admin_usecase_mock.go -package=mocks IAdminUsecase
type IAdminUsecase interface {
    GetPendingProducts(ctx context.Context, offset int) (dto.ProductsResponse, error)
    UpdateProductStatus(ctx context.Context, req dto.UpdateProductStatusRequest) error
    GetPendingUsers(ctx context.Context, offset int) (dto.UsersResponse, error)
    UpdateUserRole(ctx context.Context, req dto.UpdateUserRoleRequest) error
}

type AdminService struct {
    uc IAdminUsecase
}

func NewAdminService(uc IAdminUsecase) *AdminService {
    return &AdminService{uc: uc}
}

func (h *AdminService) GetPendingProducts(w http.ResponseWriter, r *http.Request) {
    const op = "AdminService.GetPendingProducts"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    vars := mux.Vars(r)
    offsetStr := vars["offset"]
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
        response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
        return
    }

    products, err := h.uc.GetPendingProducts(r.Context(), offset)
    if err != nil {
        logger.WithError(err).Error("get pending products")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    response.SendJSONResponse(r.Context(), w, http.StatusOK, products)
}

func (h *AdminService) UpdateProductStatus(w http.ResponseWriter, r *http.Request) {
    const op = "AdminService.UpdateProductStatus"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    var req dto.UpdateProductStatusRequest
    if err := request.ParseData(r, &req); err != nil {
        logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
        return
    }

    if err := h.uc.UpdateProductStatus(r.Context(), req); err != nil {
        logger.WithError(err).Error("update product status")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}

func (h *AdminService) GetPendingUsers(w http.ResponseWriter, r *http.Request) {
    const op = "AdminService.GetPendingUsers"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    vars := mux.Vars(r)
    offsetStr := vars["offset"]
    offset, err := strconv.Atoi(offsetStr)
    if err != nil {
        logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
        response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
        return
    }

    users, err := h.uc.GetPendingUsers(r.Context(), offset)
    if err != nil {
        logger.WithError(err).Error("get pending users")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    response.SendJSONResponse(r.Context(), w, http.StatusOK, users)
}

func (h *AdminService) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
    const op = "AdminService.UpdateUserRole"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    var req dto.UpdateUserRoleRequest
    if err := request.ParseData(r, &req); err != nil {
        logger.WithError(err).Error("parse request data")
        response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
        return
    }

    if err := h.uc.UpdateUserRole(r.Context(), req); err != nil {
        logger.WithError(err).Error("update user role")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    response.SendJSONResponse(r.Context(), w, http.StatusOK, nil)
}