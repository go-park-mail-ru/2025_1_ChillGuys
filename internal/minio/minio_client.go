package minio

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"sync"
	"time"
)

//go:generate mockgen -source=minio_client.go -destination=../minio/mocks/minio_client_mock.go -package=mocks Client
// Client интерфейс для взаимодействия с Minio
type Client interface {
    CreateOne(ctx context.Context, file FileDataType) (*UploadResponse, error)
    CreateMany(ctx context.Context, files map[string]FileDataType) ([]string, error)
    GetOne(ctx context.Context, objectID string) (string, error)
    GetMany(ctx context.Context, objectIDs []string) ([]string, error)
    DeleteOne(ctx context.Context, objectID string) error
    DeleteMany(ctx context.Context, objectIDs []string) error
}

// minioClient реализация интерфейса MinioClient
type minioClient struct {
	mc     *minio.Client
	config *config.MinioConfig
}

// NewMinioClient создает новый экземпляр Minio Client
func NewMinioClient(config *config.MinioConfig) (Client, error) {
	client, err := initMinio(config)
    if err != nil {
        return nil, err
    }
    return &minioClient{
        mc:     client,
        config: config,
    }, nil
}

// InitMinio подключается к Minio и создает бакет, если не существует
// Бакет - это контейнер для хранения объектов в Minio. Он представляет собой пространство имен, в котором можно хранить и организовывать файлы и папки.
func initMinio(config *config.MinioConfig) (*minio.Client, error) {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.RootUser, config.RootPassword, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверка наличия бакета и его создание, если не существует
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return client, nil
}

// Контекст используется для передачи сигналов об отмене операции загрузки в случае необходимости.

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру fileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
// Все операции выполняются в контексте задачи.
func (m *minioClient) CreateOne(ctx context.Context, file FileDataType) (*UploadResponse, error) {
	// Генерация уникального идентификатора для нового объекта.
	objectID := uuid.New().String()

	// Создание потока данных для загрузки в бакет Minio.
	reader := bytes.NewReader(file.Data)

	// Загрузка данных в бакет Minio с использованием контекста для возможности отмены операции.
	_, err := m.mc.PutObject(ctx, m.config.BucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create object %s: %v", file.FileName, err)
	}

	// Получение URL для загруженного объекта
	url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate URL for object %s: %v", file.FileName, err)
	}

	return &UploadResponse{
        URL:      url.String(),
        ObjectID: objectID,
    }, nil
}

// CreateMany создает несколько объектов в хранилище MinIO из переданных данных.
// Если происходит ошибка при создании объекта, метод возвращает ошибку,
// указывающую на неудачные объекты.
func (m *minioClient) CreateMany(ctx context.Context, data map[string]FileDataType) ([]string, error) {
	urls := make([]string, 0, len(data)) // Массив для хранения URL-адресов

	ctx, cancel := context.WithCancel(ctx) // Создание контекста с возможностью отмены операции.
	defer cancel()                         // Отложенный вызов функции отмены контекста при завершении функции CreateMany.

	// Создание канала для передачи URL-адресов с размером, равным количеству переданных данных.
	urlCh := make(chan string, len(data))

	var wg sync.WaitGroup // WaitGroup для ожидания завершения всех горутин.

	// Запуск горутин для создания каждого объекта.
	for objectID, file := range data {
		wg.Add(1) // Увеличение счетчика WaitGroup перед запуском каждой горутины.
		go func(objectID string, file FileDataType) {
			defer wg.Done()                                                                                                                           // Уменьшение счетчика WaitGroup после завершения горутины.
			_, err := m.mc.PutObject(ctx, m.config.BucketName, objectID, bytes.NewReader(file.Data), int64(len(file.Data)), minio.PutObjectOptions{}) // Создание объекта в бакете MinIO.
			if err != nil {
				cancel() // Отмена операции при возникновении ошибки.
				return
			}

			// Получение URL для загруженного объекта
			url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
			if err != nil {
				cancel() // Отмена операции при возникновении ошибки.
				return
			}

			urlCh <- url.String() // Отправка URL-адреса в канал с URL-адресами.
		}(objectID, file) // Передача данных объекта в анонимную горутину.
	}

	// Ожидание завершения всех горутин и закрытие канала с URL-адресами.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0.
		close(urlCh) // Закрытие канала с URL-адресами после завершения всех горутин.
	}()

	// Сбор URL-адресов из канала.
	for url := range urlCh {
		urls = append(urls, url) // Добавление URL-адреса в массив URL-адресов.
	}

	return urls, nil
}

// GetOne получает один объект из бакета Minio по его идентификатору.
// Он принимает строку `objectID` в качестве параметра и возвращает срез байт данных объекта и ошибку, если такая возникает.
func (m *minioClient) GetOne(ctx context.Context, objectID string) (string, error) {
	// Получение предварительно подписанного URL для доступа к объекту Minio.
	url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL for object %s: %v", objectID, err)
	}

	return url.String(), nil
}

// GetMany получает несколько объектов из бакета Minio по их идентификаторам.
func (m *minioClient) GetMany(ctx context.Context, objectIDs []string) ([]string, error) {
	// Создание каналов для передачи URL-адресов объектов и ошибок
	urlCh := make(chan string, len(objectIDs))         // Канал для URL-адресов объектов
	errCh := make(chan OperationError, len(objectIDs)) // Канал для ошибок

	var wg sync.WaitGroup                // WaitGroup для ожидания завершения всех горутин
	_, cancel := context.WithCancel(ctx) // Создание контекста с возможностью отмены операции
	defer cancel()                       // Отложенный вызов функции отмены контекста при завершении функции GetMany

	// Запуск горутин для получения URL-адресов каждого объекта.
	for _, objectID := range objectIDs {
		wg.Add(1) // Увеличение счетчика WaitGroup перед запуском каждой горутины
		go func(objectID string) {
			defer wg.Done()                     // Уменьшение счетчика WaitGroup после завершения горутины
			url, err := m.GetOne(ctx, objectID) // Получение URL-адреса объекта по его идентификатору с помощью метода GetOne
			if err != nil {
				errCh <- OperationError{ObjectID: objectID, Error: fmt.Errorf("failed to retrieve object %s: %v", objectID, err)} // Отправка ошибки в канал с ошибками
				cancel()                                                                                                          // Отмена операции при возникновении ошибки
				return
			}
			urlCh <- url // Отправка URL-адреса объекта в канал с URL-адресами
		}(objectID) // Передача идентификатора объекта в анонимную горутину
	}

	// Закрытие каналов после завершения всех горутин.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0
		close(urlCh) // Закрытие канала с URL-адресами после завершения всех горутин
		close(errCh) // Закрытие канала с ошибками после завершения всех горутин
	}()

	// Сбор URL-адресов объектов и ошибок из каналов.
	var urls []string // Массив для хранения URL-адресов
	var errs []error  // Массив для хранения ошибок
	for url := range urlCh {
		urls = append(urls, url) // Добавление URL-адреса в массив URL-адресов
	}
	for opErr := range errCh {
		errs = append(errs, opErr.Error) // Добавление ошибки в массив ошибок
	}

	// Проверка наличия ошибок.
	if len(errs) > 0 {
		return nil, fmt.Errorf("errors occurred while retrieving objects: %v", errs) // Возврат ошибки, если возникли ошибки при получении объектов
	}

	return urls, nil // Возврат массива URL-адресов, если ошибок не возникло
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minioClient) DeleteOne(ctx context.Context, objectID string) error {
	// Удаление объекта из бакета Minio.
	err := m.mc.RemoveObject(ctx, m.config.BucketName, objectID, minio.RemoveObjectOptions{})
	if err != nil {
		return err // Возвращаем ошибку, если не удалось удалить объект.
	}
	return nil // Возвращаем nil, если объект успешно удалён.
}

// DeleteMany удаляет несколько объектов из бакета Minio по их идентификаторам с использованием горутин.
func (m *minioClient) DeleteMany(ctx context.Context, objectIDs []string) error {
	// Создание канала для передачи ошибок с размером, равным количеству объектов для удаления
	errCh := make(chan OperationError, len(objectIDs)) // Канал для ошибок
	var wg sync.WaitGroup                              // WaitGroup для ожидания завершения всех горутин

	ctx, cancel := context.WithCancel(ctx) // Создание контекста с возможностью отмены операции
	defer cancel()                         // Отложенный вызов функции отмены контекста при завершении функции DeleteMany

	// Запуск горутин для удаления каждого объекта.
	for _, objectID := range objectIDs {
		wg.Add(1) // Увеличение счетчика WaitGroup перед запуском каждой горутины
		go func(id string) {
			defer wg.Done()                                                                     // Уменьшение счетчика WaitGroup после завершения горутины
			err := m.mc.RemoveObject(ctx, m.config.BucketName, id, minio.RemoveObjectOptions{}) // Удаление объекта с использованием Minio клиента
			if err != nil {
				errCh <- OperationError{ObjectID: id, Error: fmt.Errorf("failed to delete object %s: %v", id, err)} // Отправка ошибки в канал с ошибками
				cancel()                                                                                            // Отмена операции при возникновении ошибки
			}
		}(objectID) // Передача идентификатора объекта в анонимную горутину
	}

	// Ожидание завершения всех горутин и закрытие канала с ошибками.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0
		close(errCh) // Закрытие канала с ошибками после завершения всех горутин
	}()

	// Сбор ошибок из канала.
	var errs []error // Массив для хранения ошибок
	for opErr := range errCh {
		errs = append(errs, opErr.Error) // Добавление ошибки в массив ошибок
	}

	// Проверка наличия ошибок.
	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while deleting objects: %v", errs) // Возврат ошибки, если возникли ошибки при удалении объектов
	}

	return nil // Возврат nil, если ошибок не возникло
}

type FileDataType struct {
	FileName string
	Data     []byte
}

type OperationError struct {
	ObjectID string
	Error    error
}

type UploadResponse struct {
    URL      string `json:"url"`
    ObjectID string `json:"object_id"`
}