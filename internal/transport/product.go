package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
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
		http.Error(w, "Failed get all products", http.StatusInternalServerError)
		return
	}

	response := models.ConvertToProductsResponse(products)

	resp, err := json.Marshal(response)
	if err != nil {
		h.log.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := h.Repo.GetProductByID(r.Context(), id)
	if err != nil {
		h.log.Warnf("Product not found (ID: %d): %v", id, err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	productJson, err := json.Marshal(product)
	if err != nil{
		h.log.Errorf("Failed to encode product (ID: %d): %v", id, err)
		http.Error(w, "Failed to encode", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(productJson)
}

func (h *ProductHandler) GetProductCover(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
    idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warnf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	fileData, err := h.Repo.GetProductCoverPath(r.Context(), id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			h.log.Errorf("Cover file not found (ID: %d): %v", id, err)
			http.Error(w, "Cover file not found", http.StatusNotFound)
		} else {
			h.log.Errorf("Failed to get cover file (ID: %d): %v", id, err)
			http.Error(w, "Failed to get cover file", http.StatusInternalServerError)
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