package product

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport/utils/request"

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
	GetAllProducts(ctx context.Context, offset int) ([]*models.Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*models.Product, error)
	GetProductsByCategory(ctx context.Context, id uuid.UUID, offset int) ([]*models.Product, error)
	GetProductsByIDs(ctx context.Context, ids []uuid.UUID) ([]*models.Product, error)
	GetProductsByCategoryWithFilterAndSort(
		ctx context.Context, 
		id uuid.UUID, 
		offset int,
		minPrice, maxPrice float64,
		minRating float32,
		sortOption models.SortOption,
	) ([]*models.Product, error)
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
//	@Description	Возвращает список всех доступных продуктов
//	@Tags			products
//	@Produce		json
//	@Success		200	{array}		models.Product
//	@Failure		500	{object}	object
//	@Router			/products [get]
func (h *ProductService) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	const op = "ProductService.GetAllProducts"
	logger := logctx.GetLogger(r.Context()).WithField("op", op)

	vars := mux.Vars(r)
	offsetStr := vars["offset"]
	offset := 0
	var err error
    if offsetStr != "" {
        offset, err = strconv.Atoi(offsetStr)
        if err != nil {
            logger.WithError(err).WithField("offset", offsetStr).Error("parse offset")
            response.HandleDomainError(r.Context(), w, errs.ErrParseRequestData, op)
            return
        }
    }

	products, err := h.u.GetAllProducts(r.Context(), offset)
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
//	@Description	Возвращает детальную информацию о продукте по его ID
//	@Tags			products
//	@Produce		json
//	@Param			id	path		string	true	"UUID продукта"
//	@Success		200	{object}	models.Product
//	@Failure		400	{object}	object	"Некорректный формат UUID"
//	@Failure		404	{object}	object	"Продукт не найден"
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
//	@Description	Возвращает список товаров указанной категории, отсортированных по дате обновления (новые сначала)
//	@Tags			products
//	@Produce		json
//	@Param			id	path		string	true	"UUID категории"
//	@Success		200	{array}		models.Product
//	@Failure		400	{object}	object	"Некорректный формат UUID"
//	@Failure		404	{object}	object	"Категория не найдена"
//	@Failure		500	{object}	object
//	@Router			/products/category/{id} [get]
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

	products, err := h.u.GetProductsByCategory(r.Context(), id, offset)
	if err != nil {
		logger.WithError(err).Error("get products by category")
		response.HandleDomainError(r.Context(), w, err, op)
		return
	}

	productResponse := dto.ConvertToProductsResponse(products)

	response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}

// CreateOne godoc
//
//	@Summary		Загрузить изображение товара
//	@Description	Загружает изображение товара в хранилище MinIO
//	@Tags			products
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file				true	"Изображение товара"
//	@Success		200		{object}	map[string]string	"URL загруженного изображения"
//	@Failure		400		{object}	object				"Ошибка в данных запроса"
//	@Failure		500		{object}	object
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

// GetProductsByIDs godoc
//
//	@Summary		Получить товары по списку ID
//	@Description	Возвращает список товаров по переданным идентификаторам
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			request			body		dto.GetProductsByIDRequest	true	"Список ID товаров"
//	@Param			X-Csrf-Token	header		string						true	"CSRF-токен для защиты от подделки запросов"
//	@Success		200				{array}		dto.ProductsResponse
//	@Failure		400				{object}	object	"Некорректные данные"
//	@Failure		500				{object}	object
//	@Router			/products/batch [post]
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


func (h *ProductService) GetProductsByCategoryWithFilterAndSort(w http.ResponseWriter, r *http.Request) {
    const op = "ProductService.GetProductsByCategoryWithFilterAndSort"
    logger := logctx.GetLogger(r.Context()).WithField("op", op)

    // Парсинг ID категории
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := uuid.Parse(idStr)
    if err != nil {
        logger.WithError(err).WithField("category_id", idStr).Error("parse category ID")
        response.HandleDomainError(r.Context(), w, errs.ErrInvalidID, op)
        return
    }

    // Парсинг offset
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

    // Парсинг фильтров
    minPrice, _ := strconv.ParseFloat(r.URL.Query().Get("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(r.URL.Query().Get("max_price"), 64)
    minRating, _ := strconv.ParseFloat(r.URL.Query().Get("min_rating"), 32)

    // Парсинг параметра сортировки
    sortOption := models.SortOption(r.URL.Query().Get("sort"))
    switch sortOption {
    case models.SortByPriceAsc, models.SortByPriceDesc, models.SortByRatingAsc, models.SortByRatingDesc, models.SortByDefault:
        // допустимые значения
    default:
        sortOption = models.SortByDefault
    }

    products, err := h.u.GetProductsByCategoryWithFilterAndSort(
        r.Context(), 
        id, 
        offset,
        minPrice,
        maxPrice,
        float32(minRating),
        sortOption,
    )
    if err != nil {
        logger.WithError(err).Error("get products by category with filter and sort")
        response.HandleDomainError(r.Context(), w, err, op)
        return
    }

    productResponse := dto.ConvertToProductsResponse(products)
    response.SendJSONResponse(r.Context(), w, http.StatusOK, productResponse)
}