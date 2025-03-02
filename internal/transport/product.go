package transport

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository"
)

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

	productsJson, err := json.Marshal(products)
	if err != nil{
		http.Error(w, "Failed to encode", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(productsJson)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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