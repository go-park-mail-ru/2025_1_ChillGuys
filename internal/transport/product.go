package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
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

type ProductHandler struct {
	Repo *repository.ProductRepo
}

func NewProductHandler(repo *repository.ProductRepo) *ProductHandler {
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
