package transport

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
)

//go:generate mockgen -source=product.go -destination=../usecase/mocks/product_usecase_mock.go -package=mocks IProductUsecase
type IProductUsecase interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductCover(ctx context.Context, id uuid.UUID) ([]byte, error)
}

type ProductHandler struct {
	u IProductUsecase
	log  *logrus.Logger
}

func NewProductHandler(u IProductUsecase, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		u: u,
		log:  log,
	}
}

// GetAllProducts godoc
//	@Summary		Получить все продукты
//	@Description	Возвращает список всех продуктов
//	@Tags			products
//	@Produce		json
//	@Success		200	{object}	[]models.Product	"Список продуктов"
//	@Failure		500	{object}	utils.ErrorResponse	"Ошибка сервера"
//	@Router			/products/ [get]
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.u.GetAllProducts(r.Context())
	if err != nil {
		h.log.Warnf("Failed to get all products: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed get all products")
		return
	}

	response := models.ConvertToProductsResponse(products)

	utils.SendSuccessResponse(w, http.StatusOK, response)
}

// GetProductByID godoc
//	@Summary		Получить продукт по ID
//	@Description	Возвращает продукт по его ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path		int					true	"ID продукта"
//	@Success		200	{object}	models.Product		"Информация о продукте"
//	@Failure		400	{object}	utils.ErrorResponse	"Некорректный ID"
//	@Failure		404	{object}	utils.ErrorResponse	"Продукт не найден"
//	@Router			/products/{id} [get]
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	product, err := h.u.GetProductByID(r.Context(), id)
	if err != nil {
		h.log.Warnf("Product not found (ID: %d): %v", id, err)
		utils.SendErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	utils.SendSuccessResponse(w, http.StatusOK, product)
}

// GetProductCover godoc
//	@Summary		Получить обложку продукта
//	@Description	Возвращает обложку продукта по его ID
//	@Tags			products
//	@Produce		image/jpeg
//	@Param			id	path		int					true	"ID продукта"
//	@Success		200	{file}		[]byte				"Обложка продукта"
//	@Failure		400	{object}	utils.ErrorResponse	"Некорректный ID"
//	@Failure		404	{object}	utils.ErrorResponse	"Обложка не найдена"
//	@Failure		500	{object}	utils.ErrorResponse	"Ошибка сервера"
//	@Router			/products/{id}/cover [get]
func (h *ProductHandler) GetProductCover(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	fileData, err := h.u.GetProductCover(r.Context(), id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.log.Errorf("Cover file not found (ID: %d): %v", id, err)
			utils.SendErrorResponse(w, http.StatusNotFound, "Cover file not found")
			return
		}

		h.log.Errorf("Failed to get cover file (ID: %d): %v", id, err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get cover file")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	// Копируем содержимое файла в ответ
	if _, err := w.Write(fileData); err != nil {
		h.log.Errorf("Failed to send cover file (ID: %d): %v", id, err)
		http.Error(w, "Failed to send cover file", http.StatusInternalServerError)
		return
	}
}
