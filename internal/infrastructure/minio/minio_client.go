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
type Client interface {
	CreateOne(context.Context, FileDataType) (*UploadResponse, error)
	CreateMany(context.Context, map[string]FileDataType) ([]string, error)
	GetOne(context.Context, string) (string, error)
	GetMany(context.Context, []string) ([]string, error)
	DeleteOne(context.Context, string) error
	DeleteMany(context.Context, []string) error
}

type minioClient struct {
	mc     *minio.Client
	config *config.MinioConfig
}

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
	_, err := m.mc.PutObject(ctx, m.config.BucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
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
	urls := make([]string, 0, len(data))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Создание канала для передачи URL-адресов с размером, равным количеству переданных данных.
	urlCh := make(chan string, len(data))

	var wg sync.WaitGroup

	// Запуск горутин для создания каждого объекта.
	for objectID, file := range data {
		wg.Add(1)
		go func(objectID string, file FileDataType) {
			defer wg.Done() // Уменьшение счетчика WaitGroup после завершения горутины.
			_, err := m.mc.PutObject(ctx, m.config.BucketName, objectID, bytes.NewReader(file.Data), int64(len(file.Data)), minio.PutObjectOptions{
				ContentType: "image/jpeg",
			}) // Создание объекта в бакете MinIO.
			if err != nil {
				cancel()
				return
			}

			// Получение URL для загруженного объекта
			url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
			if err != nil {
				cancel()
				return
			}

			urlCh <- url.String()
		}(objectID, file)
	}

	// Ожидание завершения всех горутин и закрытие канала с URL-адресами.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0.
		close(urlCh) // Закрытие канала с URL-адресами после завершения всех горутин.
	}()

	for url := range urlCh {
		urls = append(urls, url)
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
	urlCh := make(chan string, len(objectIDs))
	errCh := make(chan OperationError, len(objectIDs))

	var wg sync.WaitGroup                // WaitGroup для ожидания завершения всех горутин
	_, cancel := context.WithCancel(ctx) // Создание контекста с возможностью отмены операции
	defer cancel()                       // Отложенный вызов функции отмены контекста при завершении функции GetMany

	// Запуск горутин для получения URL-адресов каждого объекта.
	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID string) {
			defer wg.Done()
			url, err := m.GetOne(ctx, objectID)
			if err != nil {
				errCh <- OperationError{ObjectID: objectID, Error: fmt.Errorf("failed to retrieve object %s: %v", objectID, err)}
				cancel()
				return
			}
			urlCh <- url
		}(objectID)
	}

	// Закрытие каналов после завершения всех горутин.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0
		close(urlCh) // Закрытие канала с URL-адресами после завершения всех горутин
		close(errCh) // Закрытие канала с ошибками после завершения всех горутин
	}()

	// Сбор URL-адресов объектов и ошибок из каналов.
	var urls []string
	var errs []error
	for url := range urlCh {
		urls = append(urls, url)
	}
	for opErr := range errCh {
		errs = append(errs, opErr.Error)
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("errors occurred while retrieving objects: %v", errs)
	}

	return urls, nil
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minioClient) DeleteOne(ctx context.Context, objectID string) error {
	// Удаление объекта из бакета Minio.
	if err := m.mc.RemoveObject(ctx, m.config.BucketName, objectID, minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

// DeleteMany удаляет несколько объектов из бакета Minio по их идентификаторам с использованием горутин.
func (m *minioClient) DeleteMany(ctx context.Context, objectIDs []string) error {
	// Создание канала для передачи ошибок с размером, равным количеству объектов для удаления
	errCh := make(chan OperationError, len(objectIDs))
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Запуск горутин для удаления каждого объекта.
	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if err := m.mc.RemoveObject(ctx, m.config.BucketName, id, minio.RemoveObjectOptions{}); err != nil {
				errCh <- OperationError{ObjectID: id, Error: fmt.Errorf("failed to delete object %s: %v", id, err)}
				cancel()
			}
		}(objectID)
	}

	// Ожидание завершения всех горутин и закрытие канала с ошибками.
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Сбор ошибок из канала.
	var errs []error
	for opErr := range errCh {
		errs = append(errs, opErr.Error)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while deleting objects: %v", errs)
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
