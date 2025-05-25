package transport

import (
	"context"
	"github.com/mailru/easyjson"
	"io"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/usecase/helpers"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

//go:generate mockgen -source=seller.go -destination=../../usecase/mocks/seller_usecase_mock.go -package=mocks ISellerUsecase
type ISellerUsecase interface {
	AddProduct(ctx context.Context, product *models.Product, categoryID uuid.UUID) (*models.Product, error)
	UploadProductImage(ctx context.Context, productID uuid.UUID, imageURL string) error
	GetSellerProducts(ctx context.Context, sellerID uuid.UUID, offset int) ([]*models.Product, error)
	CheckProductBelongs(ctx context.Context, productID, sellerID uuid.UUID) (bool, error)
}

type SellerHandler struct {
	usecase      ISellerUsecase
	minioService minio.Provider
}

func NewSellerHandler(u ISellerUsecase, ms minio.Provider) *SellerHandler {
	return &SellerHandler{
		usecase:      u,
		minioService: ms,
	}
}

// AddProduct godoc
// @Summary Добавить товар (без изображения)
// @Description Добавляет новый товар без изображения
// @Tags seller
// @Accept json
// @Produce json
// @Param product body dto.AddProductRequest true "Данные товара"
// @Success 201 {object} models.Product
// @Failure 400 {object} object
// @Failure 403 {object} object
// @Failure 500 {object} object
// @Router /seller/products [post]
func (h *SellerHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	const op = "SellerHandler.AddProduct"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	// Получаем ID продавца из контекста
	sellerID, err := helpers.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		response.HandleDomainError(r.Context(), w, err, op)
	}

	var req dto.AddProductRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &req); err != nil {
		logger.WithError(err).Error("failed to parse request data")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	// Конвертируем в модель
	product := &models.Product{
		SellerID:        sellerID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Quantity:        req.Quantity,
		PreviewImageURL: "", // Будет добавлено позже
	}

	categoryID, err := uuid.Parse(req.Category)
	if err != nil {
		logger.WithError(err).Error("parse category ID")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	newProduct, err := h.usecase.AddProduct(r.Context(), product, categoryID)
	if err != nil {
		logger.WithError(err).Error("add product")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusCreated, newProduct)
}

// UploadProductImage godoc
// @Summary Загрузить изображение товара
// @Description Загружает изображение для указанного товара
// @Tags seller
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID товара"
// @Param file formData file true "Изображение товара"
// @Success 200 {object} map[string]string
// @Failure 400 {object} object
// @Failure 403 {object} object
// @Failure 500 {object} object
// @Router /seller/products/{id}/image [post]
func (h *SellerHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	const op = "SellerHandler.UploadProductImage"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	// Получаем ID продавца из контекста
	sellerID, err := helpers.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		response.HandleDomainError(r.Context(), w, err, op)
	}

	// Получаем ID товара из URL
	vars := mux.Vars(r)
	productID, err := uuid.Parse(vars["id"])
	if err != nil {
		logger.WithError(err).Error("parse product ID")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Проверяем, что товар принадлежит продавцу
	belongs, err := h.usecase.CheckProductBelongs(r.Context(), productID, sellerID)
	if err != nil || !belongs {
		logger.WithError(err).Error("product doesn't belong to seller")
		response.SendJSONError(r.Context(), w, http.StatusForbidden, "product doesn't belong to seller")
		return
	}

	// Обработка загрузки файла
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		logger.WithError(err).Error("parse multipart form")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		logger.WithError(err).Error("get file from form")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "no file uploaded")
		return
	}
	defer file.Close()

	// Читаем файл и загружаем
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("read file content")
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "failed to read file")
		return
	}

	fileData := minio.FileData{
		Name: header.Filename,
		Data: fileBytes,
	}

	// Загружаем файл в MinIO
	productResponse, err := h.minioService.CreateOne(r.Context(), fileData)
	if err != nil {
		logger.WithError(err).Error("upload file to minio")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	err = h.usecase.UploadProductImage(r.Context(), productID, productResponse.URL)
	if err != nil {
		logger.WithError(err).Error("upload product image")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}

// GetSellerProducts godoc
// @Summary Получить товары продавца
// @Description Возвращает список всех товаров текущего продавца
// @Tags seller
// @Produce json
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {array} models.Product
// @Failure 403 {object} object
// @Failure 500 {object} object
// @Router /seller/products [get]
func (h *SellerHandler) GetSellerProducts(w http.ResponseWriter, r *http.Request) {
	const op = "SellerHandler.GetSellerProducts"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	// Получаем ID продавца из контекста
	sellerID, err := helpers.GetUserIDFromContext(r.Context())
	if err != nil {
		logger.WithError(err).Error("get user ID from context")
		response.HandleDomainError(r.Context(), w, err, op)
	}

	vars := mux.Vars(r)
	offsetStr := vars["offset"]
	offset := 0
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
			response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
			return
		}
	}

	products, err := h.usecase.GetSellerProducts(r.Context(), sellerID, offset)
	if err != nil {
		logger.WithError(err).Error("get seller products")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	productResponse := dto.ConvertToSellerProductsResponse(products)
	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}
