package category

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
)

//go:generate mockgen -source=category.go -destination=../../usecase/mocks/category_usecase_mock.go -package=mocks ICategoryUsecase
type ICategoryUsecase interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
}

type CategoryService struct {
	u            ICategoryUsecase
}

func NewCategoryService(u ICategoryUsecase) *CategoryService {
	return &CategoryService{
		u:            u,
	}
}

// GetAllCategories godoc
//
//	@Summary		Получить все категории
//	@Description	Возвращает список всех доступных категорий товаров
//	@Tags			categories
//	@Produce		json
//	@Success		200	{object}	[]models.Category	"Список категорий"
//	@Failure		500	{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
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