package product

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"
	"io"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/infrastructure/minio"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models/errs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/dto"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/middleware/logctx"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/response"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
)

//go:generate mockgen -source=product.go -destination=../../usecase/mocks/product_usecase_mock.go -package=mocks IProductUsecase
type IProductUsecase interface {
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductsByCategory(ctx context.Context, id uuid.UUID) ([]*models.Product, error)
	GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error)
}

type ProductService struct {
	u            IProductUsecase
	minioService minio.Provider
}

func NewProductService(u IProductUsecase, ms minio.Provider) *ProductService {
	return &ProductService{
		u:            u,
		minioService: ms,
	}
}

// GetAllProducts godoc
//
//	@Summary		Получить все продукты
//	@Description	Возвращает список всех продуктов
//	@Tags			products
//	@Produce		json
//	@Success		200	{object}	[]models.Product	"Список продуктов"
//	@Failure		500	{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/products [get]
func (h *ProductService) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.GetAllProducts"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	products, err := h.u.GetAllProducts(r.Context())
	if err != nil {
		logger.WithError(err).Error("get all products")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	productResponse := dto.ConvertToProductsResponse(products)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}

// GetProductByID godoc
//
//	@Summary		Получить продукт по ID
//	@Description	Возвращает продукт по его ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path		string				true	"ID продукта"
//	@Success		200	{object}	models.Product		"Информация о продукте"
//	@Failure		400	{object}	dto.ErrorResponse	"Некорректный ID"
//	@Failure		404	{object}	dto.ErrorResponse	"Продукт не найден"
//	@Router			/products/{id} [get]
func (h *ProductService) GetProductByID(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.GetProductByID"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithError(err).WithField("product_id", idStr).Error("parse product ID")
		response.HandleDomainError(r.Context(), w, errs.ErrInvalidID, op)
		return
	}

	logger = logger.WithField("product_id", id)
	product, err := h.u.GetProductByID(r.Context(), id)
	if err != nil {
		logger.WithError(err).Error("get product by ID")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, product)
}

// GetProductsByCategory godoc
//
//	@Summary		Получить товары по категории
//	@Description	Возвращает список всех одобренных товаров, принадлежащих указанной категории.
//
// Товары сортируются по дате обновления (сначала новые).
//
//	@Tags			products
//	@Produce		json
//	@Param			id	path		string				true	"UUID категории в формате строки"
//	@Success		200	{object}	[]models.Product	"Успешный запрос. Возвращает массив товаров."
//	@Failure		400	{object}	dto.ErrorResponse	"Неверный формат UUID категории"
//	@Failure		404	{object}	dto.ErrorResponse	"Категория не найдена"
//	@Failure		500	{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/v1/products/category/{id} [get]
func (h *ProductService) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.GetProductsByCategory"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.WithError(err).WithField("category_id", idStr).Error("parse category ID")
		response.HandleDomainError(r.Context(), w, errs.ErrInvalidID, op)
		return
	}

	products, err := h.u.GetProductsByCategory(r.Context(), id)
	if err != nil {
		logger.WithError(err).Error("get products by category")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	productResponse := dto.ConvertToProductsResponse(products)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}

// FIXME: models.SuccessResponse не найден

// CreateOne godoc
//
//	@Summary		Загрузить файл в MinIO
//	@Description	Загружает один файл в хранилище MinIO
//	@Tags			products
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file				true	"Файл для загрузки"
//	@Success		200		{object}	map[string]string	"Информация о загруженном файле"
//	@Failure		400		{object}	dto.ErrorResponse	"Ошибка в запросе"
//	@Failure		500		{object}	dto.ErrorResponse	"Ошибка сервера"
//	@Router			/products/upload [post]
func (h *ProductService) CreateOne(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.CreateOne"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	// Проверяем, что запрос содержит multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.WithError(err).Error("parse multipart form")
		response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.WithError(err).Error("get file from form")
		response.HandleDomainError(r.Context(), w, fmt.Errorf("no file uploaded"), op)
		return
	}
	defer file.Close()

	// Читаем содержимое файла
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logger.WithError(err).Error("read file content")
		response.HandleDomainError(r.Context(), w, fmt.Errorf("failed to read file"), op)
		return
	}

	// Создаем структуру для MinIO
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

	// Возвращаем успешный ответ с URL файла
	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}

func (p *ProductService) GetProductsByIDs(w http.ResponseWriter, r *http.Request) {
	var req dto.GetProductsByIDRequest
	if err := request.ParseData(r, &req); err != nil {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.ProductIDs) == 0 {
		response.SendJSONError(r.Context(), w, http.StatusBadRequest, "at least one product ID is required")
		return
	}

	products, err := p.u.GetProductsByIDs(r.Context(), req.ProductIDs)
	if err != nil {
		response.HandleDomainError(r.Context(), w, err, "failed to get products by IDs")
		return
	}

	response.SendJSONResponse(r.Context(), w, http.StatusOK, dto.ConvertToProductsResponse(products))
}
