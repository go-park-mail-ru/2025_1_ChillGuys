package tests

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/models"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/repository/mocks"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/transport"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var testProducts = []*models.Product{
	{ID: 1, Name: "Смартфон Xiaomi Redmi Note 10", Description: "Смартфон с AMOLED-дисплеем и камерой 48 Мп", Count: 50, Price: 19999, ReviewsCount: 120, Rating: 4.5},
	{ID: 2, Name: "Ноутбук ASUS VivoBook 15", Description: "Ноутбук с процессором Intel Core i5 и SSD на 512 ГБ", Count: 30, Price: 54999, ReviewsCount: 80, Rating: 4.7},
	{ID: 3, Name: "Наушники Sony WH-1000XM4", Description: "Беспроводные наушники с шумоподавлением", Count: 25, Price: 29999, ReviewsCount: 200, Rating: 4.8},
	{ID: 4, Name: "Фитнес-браслет Xiaomi Mi Band 6", Description: "Фитнес-браслет с AMOLED-дисплеем и мониторингом сна", Count: 100, Price: 3999, ReviewsCount: 300, Rating: 4.6},
	{ID: 5, Name: "Пылесос Dyson V11", Description: "Беспроводной пылесос с мощным всасыванием", Count: 15, Price: 59999, ReviewsCount: 90, Rating: 4.9},
	{ID: 6, Name: "Кофемашина DeLonghi Magnifica", Description: "Автоматическая кофемашина для приготовления эспрессо", Count: 10, Price: 79999, ReviewsCount: 70, Rating: 4.7},
	{ID: 7, Name: "Электросамокат Xiaomi Mi Scooter 3", Description: "Электросамокат с запасом хода 30 км", Count: 40, Price: 29999, ReviewsCount: 150, Rating: 4.5},
	{ID: 8, Name: "Умная колонка Яндекс.Станция Мини", Description: "Умная колонка с голосовым помощником Алисой", Count: 60, Price: 7999, ReviewsCount: 250, Rating: 4.4},
	{ID: 9, Name: "Монитор Samsung Odyssey G5", Description: "Игровой монитор с разрешением 1440p и частотой 144 Гц", Count: 20, Price: 34999, ReviewsCount: 100, Rating: 4.6},
	{ID: 10, Name: "Электрочайник Bosch TWK 3A011", Description: "Электрочайник с мощностью 2400 Вт", Count: 50, Price: 1999, ReviewsCount: 180, Rating: 4.3},
	{ID: 11, Name: "Робот-пылесос iRobot Roomba 981", Description: "Робот-пылесос с навигацией по карте помещения", Count: 12, Price: 69999, ReviewsCount: 60, Rating: 4.8},
	{ID: 12, Name: "Фен Dyson Supersonic", Description: "Фен с технологией защиты волос от перегрева", Count: 18, Price: 49999, ReviewsCount: 130, Rating: 4.7},
	{ID: 13, Name: "Микроволновая печь LG MS-2042DB", Description: "Микроволновка с объемом 20 литров", Count: 35, Price: 8999, ReviewsCount: 110, Rating: 4.2},
	{ID: 14, Name: "Игровая консоль PlayStation 5", Description: "Игровая консоль нового поколения", Count: 5, Price: 79999, ReviewsCount: 300, Rating: 4.9},
	{ID: 15, Name: "Электронная книга PocketBook 740", Description: "Электронная книга с экраном E Ink Carta", Count: 25, Price: 19999, ReviewsCount: 90, Rating: 4.4},
}

func TestGetAllProducts(t *testing.T) {
	// Инициализация контроллера Gomock.
	ctrl := gomock.NewController(t)
	// Убедимся, что все ожидаемые вызовы мока были выполнены.
	defer ctrl.Finish()

	// Создание мока репозитория.
	mockRepo := mocks.NewMockIProductRepo(ctrl)
	logger := logrus.New()

	// Создание обработчика с моком репозитория.
	handler := transport.NewProductHandler(mockRepo, logger)

	t.Run("Success case", func(t *testing.T) {
				// Создание HTTP-запроса для эндпоинта /products/.
		// Метод GET, тело запроса отсутствует.
		req, err := http.NewRequest("GET", "/products/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Настройка ожидаемого поведения мока.
		// Ожидаем, что метод GetAllProducts будет вызван один раз с любым контекстом (gomock.Any()).
		// В ответ мок вернет тестовые данные (testProducts) и nil в качестве ошибки.
		mockRepo.EXPECT().GetAllProducts(req.Context()).Return(testProducts, nil).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetAllProducts(rr, req)

		// Проверка статус-кода ответа.
		// Ожидаем, что статус-код будет 200 (OK).
		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200, got %d", rr.Code)

		// Декодирование JSON-ответа.
		var response models.ProductsResponse
		err = json.Unmarshal(rr.Body.Bytes(), &response)
		// Проверка, что декодирование прошло без ошибок.
		assert.NoError(t, err, "Failed to decode response")

		// Проверка, что список продуктов не пуст.
		assert.NotEmpty(t, response.Products, "Expected non-empty product list")
		assert.Equal(t, len(testProducts), response.Total, "Expected %d products, got %d", len(testProducts), response.Total)
	})

	t.Run("Repository error case", func(t *testing.T) {
		// Создание HTTP-запроса для эндпоинта /products/.
		req, err := http.NewRequest("GET", "/products/", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Настройка ожидаемого поведения мока.
		// Ожидаем, что метод GetAllProducts будет вызван один раз с любым контекстом (gomock.Any()).
		// В ответ мок вернет ошибку.
		mockRepo.EXPECT().GetAllProducts(req.Context()).Return(nil, errors.New("repository error")).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetAllProducts(rr, req)

		// Проверка статус-кода ответа.
		// Ожидаем, что статус-код будет 500 (Internal Server Error).
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status code 500, got %d", rr.Code)

		// Проверка тела ответа.
		// Ожидаем, что в ответе будет сообщение об ошибке.
		assert.Contains(t, rr.Body.String(), "Failed get all products", "Expected error message in response")
	})
}

func TestGetProductByID(t *testing.T) {
	// Инициализация контроллера Gomock.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока репозитория.
	mockRepo := mocks.NewMockIProductRepo(ctrl)
	logger := logrus.New()

	// Создание обработчика с моком репозитория и стандартным JSONEncoder.
	handler := transport.NewProductHandler(mockRepo, logger)

	t.Run("Success case", func(t *testing.T) {
		// Тестовый ID продукта.
		testID := 1
		// Тестовый продукт.
		testProduct := &models.Product{
			ID:          testID,
			Name:        "Смартфон Xiaomi Redmi Note 10",
			Description: "Смартфон с AMOLED-дисплеем и камерой 48 Мп",
			Count:       50,
			Price:       19999,
			ReviewsCount: 120,
			Rating:      4.5,
		}

		// Создание HTTP-запроса с path-параметром.
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": strconv.Itoa(testID),
		}
		req = mux.SetURLVars(req, vars)

		// Настройка ожидаемого поведения мока.
		// Используем req.Context() вместо context.Background().
		mockRepo.EXPECT().GetProductByID(req.Context(), testID).Return(testProduct, nil).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductByID(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200, got %d", rr.Code)

		// Декодирование JSON-ответа.	
		var product models.Product
		err = json.Unmarshal(rr.Body.Bytes(), &product)
		assert.NoError(t, err, "Failed to decode response")

		// Проверка данных продукта.
		assert.Equal(t, testProduct.ID, product.ID, "Expected product ID %d, got %d", testID, product.ID)
		assert.Equal(t, testProduct.Name, product.Name, "Product name mismatch")
		assert.Equal(t, testProduct.Description, product.Description, "Product description mismatch")
	})

	t.Run("Invalid ID case", func(t *testing.T) {
		// Невалидный ID (не число).
		invalidID := "abc"

		// Создание HTTP-запроса с невалидным path-параметром.
		req, err := http.NewRequest("GET", "/products/"+invalidID, nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": invalidID,
		}
		req = mux.SetURLVars(req, vars)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductByID(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400, got %d", rr.Code)

		// Проверка тела ответа.
		assert.Contains(t, rr.Body.String(), "Invalid ID", "Expected error message in response")
	})

	t.Run("Repository error case", func(t *testing.T) {
		// Тестовый ID продукта.
		testID := 1

		// Создание HTTP-запроса с path-параметром.
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID), nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Создание HTTP-запроса с path-параметром.
		vars := map[string]string{
			"id": strconv.Itoa(testID),
		}
		req = mux.SetURLVars(req, vars)

		// Настройка ожидаемого поведения мока.
		// Используем req.Context() вместо context.Background().
		mockRepo.EXPECT().GetProductByID(req.Context(), testID).Return(nil, errors.New("not found")).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductByID(rr, req)

		// Получение результата ответа.
		resp := rr.Result()
		defer resp.Body.Close() // Закрываем тело ответа.

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status code 404, got %d", rr.Code)

		// Чтение тела ответа.
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err, "Failed to read response body")

		// Проверка тела ответа.
		assert.Contains(t, string(body), "Not found", "Expected error message in response")
	})
}

func TestGetrProductCover(t *testing.T) {
	// Инициализация контроллера Gomock.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создание мока репозитория.
	mockRepo := mocks.NewMockIProductRepo(ctrl)
	logger := logrus.New()

	// Создание обработчика с моком репозитория.
	handler := transport.NewProductHandler(mockRepo, logger)

	t.Run("Success case", func(t *testing.T) {
		// Тестовый ID продукта.
		testID := 1
		// Тестовые данные обложки (заглушка).
		testCoverData := []byte{0xFF, 0xD8, 0xFF} // Пример JPEG-заголовка.

		// Создание HTTP-запроса с path-параметром.
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID)+"/cover", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": strconv.Itoa(testID),
		}
		req = mux.SetURLVars(req, vars)

		// Настройка ожидаемого поведения мока.
		mockRepo.EXPECT().GetProductCoverPath(req.Context(), testID).Return(testCoverData, nil).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductCover(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200, got %d", rr.Code)

		// Проверка заголовка Content-Type.
		assert.Equal(t, "image/jpeg", rr.Header().Get("Content-Type"), "Expected Content-Type image/jpeg, got %s", rr.Header().Get("Content-Type"))

		// Проверка тела ответа.
		assert.Equal(t, testCoverData, rr.Body.Bytes(), "Expected response body to match test cover data")
	})

	t.Run("Cover not found case", func(t *testing.T) {
		// Тестовый ID продукта.
		testID := 2

		// Создание HTTP-запроса с path-параметром.
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID)+"/cover", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": strconv.Itoa(testID),
		}
		req = mux.SetURLVars(req, vars)

		// Настройка ожидаемого поведения мока.
		mockRepo.EXPECT().GetProductCoverPath(req.Context(), testID).Return(nil, os.ErrNotExist).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductCover(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusNotFound, rr.Code, "Expected status code 404, got %d", rr.Code)

		// Проверка тела ответа.
		assert.Contains(t, rr.Body.String(), "Cover file not found", "Expected error message in response")
	})

	t.Run("Invalid ID case", func(t *testing.T) {
		// Невалидный ID (не число).
		invalidID := "abc"

		// Создание HTTP-запроса с невалидным path-параметром.
		req, err := http.NewRequest("GET", "/products/"+invalidID+"/cover", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": invalidID,
		}
		req = mux.SetURLVars(req, vars)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductCover(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status code 400, got %d", rr.Code)

		// Проверка тела ответа.
		assert.Contains(t, rr.Body.String(), "Invalid ID", "Expected error message in response")
	})

	t.Run("Internal server error case", func(t *testing.T) {
		// Тестовый ID продукта.
		testID := 3

		// Создание HTTP-запроса с path-параметром.
		req, err := http.NewRequest("GET", "/products/"+strconv.Itoa(testID)+"/cover", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		// Добавление path-параметра в запрос.
		vars := map[string]string{
			"id": strconv.Itoa(testID),
		}
		req = mux.SetURLVars(req, vars)

		// Настройка ожидаемого поведения мока.
		mockRepo.EXPECT().GetProductCoverPath(req.Context(), testID).Return(nil, fmt.Errorf("internal error")).Times(1)

		// Создание ResponseRecorder для записи ответа.
		rr := httptest.NewRecorder()

		// Вызов обработчика.
		handler.GetProductCover(rr, req)

		// Проверка статус-кода ответа.
		assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status code 500, got %d", rr.Code)

		// Проверка тела ответа.
		assert.Contains(t, rr.Body.String(), "Failed to get cover file", "Expected error message in response")
	})
}