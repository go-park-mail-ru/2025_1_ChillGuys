package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/minio"
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
	minioService minio.Client
}

func NewProductHandler(u IProductUsecase, log *logrus.Logger, mS minio.Client) *ProductHandler {
	return &ProductHandler{
		u: u,
		log:  log,
		minioService: mS,
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



// CreateOne godoc
// @Summary Загрузить файл в MinIO
// @Description Загружает один файл в хранилище MinIO
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл для загрузки"
// @Success 200 {object} models.SuccessResponse "Файл успешно загружен"
// @Failure 400 {object} utils.ErrorResponse "Ошибка в запросе"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /products/upload [post]
func (h *ProductHandler) CreateOne(w http.ResponseWriter, r *http.Request) {
    // Проверяем, что запрос содержит multipart/form-data
    if err := r.ParseMultipartForm(10 << 20); err != nil { // Максимум 10MB файл
        h.log.Warnf("Error parsing multipart form: %v", err)
        utils.SendErrorResponse(w, http.StatusBadRequest, "Failed to parse form data")
        return
    }

    // Получаем файл из формы
    file, header, err := r.FormFile("file")
    if err != nil {
        h.log.Warnf("Error getting file from form: %v", err)
        utils.SendErrorResponse(w, http.StatusBadRequest, "No file uploaded")
        return
    }
    defer file.Close()

    // Читаем содержимое файла
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        h.log.Errorf("Error reading file: %v", err)
        utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to read file")
        return
    }

    // Создаем структуру для MinIO
    fileData := minio.FileDataType{
        FileName: header.Filename,
        Data:     fileBytes,
    }

    // Загружаем файл в MinIO
    response, err := h.minioService.CreateOne(r.Context(), fileData)
    if err != nil {
        h.log.Errorf("Upload error: %v", err)
        utils.SendErrorResponse(w, http.StatusInternalServerError, "Upload failed")
        return
    }

    // Возвращаем успешный ответ с URL файла
    utils.SendSuccessResponse(w, http.StatusOK, response)
}

// GetOne godoc
// @Summary Получить файл по ID
// @Description Возвращает URL для доступа к файлу в MinIO
// @Tags files
// @Produce json
// @Param objectID path string true "ID объекта в MinIO"
// @Success 200 {object} models.SuccessResponse "Успешный ответ с URL файла"
// @Failure 400 {object} utils.ErrorResponse "Неверный ID объекта"
// @Failure 404 {object} utils.ErrorResponse "Файл не найден"
// @Failure 500 {object} utils.ErrorResponse "Ошибка сервера"
// @Router /files/{objectID} [get]
func (h *ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) {
    // Получаем ID объекта из параметров URL
    vars := mux.Vars(r)
    objectID := vars["objectID"]

    // Проверяем валидность UUID
    if _, err := uuid.Parse(objectID); err != nil {
        h.log.Warnf("Invalid object ID format: %s", objectID)
        utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid object ID format")
        return
    }

    // Получаем URL файла из MinIO
    url, err := h.minioService.GetOne(r.Context(), objectID)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            h.log.Warnf("File not found (ID: %s): %v", objectID, err)
            utils.SendErrorResponse(w, http.StatusNotFound, "File not found")
            return
        }

        h.log.Errorf("Failed to get file (ID: %s): %v", objectID, err)
        utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get file")
        return
    }

    // Формируем успешный ответ
    utils.SendSuccessResponse(w, http.StatusOK, map[string]string{
        "url":       url,
        "object_id": objectID,
    })
}