package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/gorilla/mux"
)

type ProductsResponse struct {
	Total    int                   `json:"total"`
	Products []models.BriefProduct `json:"products"`
}

func convertToProductsResponse(products []*models.Product) ProductsResponse {
	briefProducts := make([]models.BriefProduct, 0, len(products))
	for _, product := range products {
		briefProduct := models.ConvertToBriefProduct(product)
		briefProducts = append(briefProducts, briefProduct)
	}

	response := ProductsResponse{
		Total: len(briefProducts),
		Products: briefProducts,
	}

	return response
}

//go:generate mockgen -source=product.go -destination=../repository/mocks/product_repo_mock.go package=mocks IProductRepo
type IProductRepo interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
}

type ProductHandler struct {
	Repo IProductRepo
}

func NewProductHandler(repo IProductRepo) *ProductHandler {
	return &ProductHandler{
		Repo: repo,
	}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.Repo.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "Failed get all products", http.StatusInternalServerError)
		return
	}

	response := convertToProductsResponse(products)

	resp, err := json.Marshal(response)
	if err != nil {
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
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := h.Repo.GetProductByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	productJson, err := json.Marshal(product)
	if err != nil{
		http.Error(w, "Failed to encode", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(productJson)
}
