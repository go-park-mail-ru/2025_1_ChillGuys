package category

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//go:generate mockgen -source=category.go -destination=../../usecase/mocks/category_usecase_mock.go -package=mocks ICategoryUsecase
type ICategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetAllSubategories(ctx context.Context, category_id uuid.UUID) ([]*models.Category, error)
	GetNameSubcategory(ctx context.Context, id uuid.UUID) (string, error)
}

type CategoryService struct {
	u ICategoryUsecase
}

func NewCategoryService(u ICategoryUsecase) *CategoryService {
	return &CategoryService{
		u: u,
	}
}

// GetAllCategories godoc
//
//	@Summary		Получить все категории
//	@Description	Возвращает список всех доступных категорий товаров
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		dto.CategoryResponse
//	@Failure		500	{object}	object
//	@Router			/categories [get]
func (h *CategoryService) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryService.GetAllCategories"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	categories, err := h.u.GetAllCategories(r.Context())
	if err != nil {
		logger.WithError(err).Error("get all categories")
		response.HandleDomainError(r.Context(), w, err, "get categories")
		return
	}

	categoryResponse := dto.ConvertToCategoriesResponse(categories)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, categoryResponse)
}

func (h *CategoryService) GetAllSubcategories(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryService.GetAllSubcategories"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithError(err).WithField("category_id", idStr).Error("parse category ID")
		response.HandleDomainError(r.Context(), w, errs.ErrInvalidID, op)
		return
	}

	categories, err := h.u.GetAllSubategories(r.Context(), id)
	if err != nil {
		logger.WithError(err).Error("get all subcategories")
		response.HandleDomainError(r.Context(), w, err, "get subcategories")
		return
	}

	categoryResponse := dto.ConvertToCategoriesResponse(categories)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, categoryResponse)
}

func (h *CategoryService) GetNameSubcategory(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryService.GetNameSubcategories"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithError(err).WithField("subcategory_id", idStr).Error("parse subcategory ID")
		response.HandleDomainError(r.Context(), w, errs.ErrInvalidID, op)
		return
	}

	name, err := h.u.GetNameSubcategory(r.Context(), id)
	if err != nil {
		logger.WithError(err).Error("get name subcategory")
		response.HandleDomainError(r.Context(), w, err, "get name subcategory")
		return
	}

	var resp dto.NameSubcategory;
	resp.Name = name

	response.SendJSONResponse(r.Context(), w, http.StatusOK, resp)
}