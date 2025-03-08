package transport

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=product.go -destination=../repository/mocks/product_repo_mock.go package=mocks IProductRepo
type IProductRepo interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	GetProductCoverPath(ctx context.Context, id int) ([]byte, error)
}

type ProductHandler struct {
	Repo IProductRepo
	log *logrus.Logger
}

func NewProductHandler(repo IProductRepo, log *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		Repo: repo,
		log: log,
	}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Repo.GetAllProducts(r.Context())
	if err != nil {
		h.log.Warnf("Failed to get all products: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed get all products")
		return
	}

	response := models.ConvertToProductsResponse(products)

	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	product, err := h.Repo.GetProductByID(r.Context(), id)
	if err != nil {
		h.log.Warnf("Product not found (ID: %d): %v", id, err)
		utils.SendErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}
	
	utils.SendSuccessResponse(w, http.StatusOK, product)
}

func (h *ProductHandler) GetProductCover(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
    idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	fileData, err := h.Repo.GetProductCoverPath(r.Context(), id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.log.Errorf("Cover file not found (ID: %d): %v", id, err)
			utils.SendErrorResponse(w, http.StatusNotFound, "Cover file not found")
		} else {
			h.log.Errorf("Failed to get cover file (ID: %d): %v", id, err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get cover file")
		}
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