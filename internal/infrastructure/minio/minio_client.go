package minio

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=minio_Provider.go -destination=../minio/mocks/minio_Provider_mock.go -package=mocks Provider
type Provider interface {
	CreateOne(context.Context, FileData) (*UploadResponse, error)
	CreateMany(context.Context, map[string]FileData) ([]string, error)
	GetOne(context.Context, string) (string, error)
	GetMany(context.Context, []string) ([]string, error)
	DeleteOne(context.Context, string) error
	DeleteMany(context.Context, []string) error
}

type minioProvider struct {
	mc     *minio.Client
	config *config.MinioConfig
	log    *logrus.Logger
}

func NewMinioProvider(config *config.MinioConfig, log *logrus.Logger) (Provider, error) {
	Provider, err := initMinio(config)
	if err != nil {
		return nil, err
	}
	return &minioProvider{
		mc:     Provider,
		config: config,
		log:    log,
	}, nil
}

// InitMinio подключается к Minio и создает бакет, если не существует
// Бакет - это контейнер для хранения объектов в Minio. Он представляет собой пространство имен, в котором можно хранить и организовывать файлы и папки.
func initMinio(config *config.MinioConfig) (*minio.Client, error) {
	// Создание контекста с возможностью отмены операции
	ctx := context.Background()

	// Подключение к Minio с использованием имени пользователя и пароля
	Provider, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.RootUser, config.RootPassword, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверка наличия бакета и его создание, если не существует
	exists, err := Provider.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, err
	}
	if !exists {
		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {"AWS": ["*"]},
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`
		if err = Provider.SetBucketPolicy(ctx, config.BucketName, fmt.Sprintf(policy, config.BucketName)); err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %v", err)
		}
	}

	return Provider, nil
}

// CreateOne создает один объект в бакете Minio.
// Метод принимает структуру fileData, которая содержит имя файла и его данные.
// В случае успешной загрузки данных в бакет, метод возвращает nil, иначе возвращает ошибку.
// Все операции выполняются в контексте задачи.
func (m *minioProvider) CreateOne(ctx context.Context, file FileData) (*UploadResponse, error) {
	// Генерация уникального идентификатора для нового объекта.
	objectID := uuid.New().String()
	logFields := logrus.Fields{
		"object_id": objectID,
		"file_name": file.Name,
		"size":      len(file.Data),
	}

	m.log.WithFields(logFields).Debug("attempting to upload file to MinIO")

	// Создание потока данных для загрузки в бакет Minio.
	reader := bytes.NewReader(file.Data)

	// Загрузка данных в бакет Minio с использованием контекста для возможности отмены операции.
	uploadInfo, err := m.mc.PutObject(ctx, m.config.BucketName, objectID, reader, int64(len(file.Data)), minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		m.log.WithFields(logFields).WithError(err).Error("failed to upload file to MinIO")
		return nil, fmt.Errorf("failed to create object %s: %v", file.Name, err)
	}

	// Добавляем информацию о загрузке в логи
	logFields["upload_info"] = uploadInfo
	m.log.WithFields(logFields).Debug("file upload details")

	// Получение URL для загруженного объекта
	url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
	if err != nil {
		m.log.WithFields(logFields).WithError(err).Error("failed to generate presigned URL")
		return nil, fmt.Errorf("failed to generate URL for object %s: %v", file.Name, err)
	}

	m.log.WithFields(logFields).Info("file successfully uploaded to MinIO")

	return &UploadResponse{
		URL:      url.String(),
		ObjectID: objectID,
	}, nil
}

// CreateMany создает несколько объектов в хранилище MinIO из переданных данных.
// Если происходит ошибка при создании объекта, метод возвращает ошибку,
// указывающую на неудачные объекты.
func (m *minioProvider) CreateMany(ctx context.Context, data map[string]FileData) ([]string, error) {
	logFields := logrus.Fields{
        "total_files": len(data),
    }

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Создание канала для передачи URL-адресов с размером, равным количеству переданных данных.
	urlCh := make(chan string, len(data))
	errCh := make(chan error, len(data))

	var wg sync.WaitGroup

	// Запуск горутин для создания каждого объекта.
	for objectID, file := range data {
		wg.Add(1)
		go func(objectID string, file FileData) {
			defer wg.Done() // Уменьшение счетчика WaitGroup после завершения горутины.

			fileLogFields := logrus.Fields{
				"object_id": objectID,
				"file_name": file.Name,
				"size":      len(file.Data),
			}

			uploadInfo, err := m.mc.PutObject(
				ctx,
				m.config.BucketName,
				objectID,
				bytes.NewReader(file.Data),
				int64(len(file.Data)),
				minio.PutObjectOptions{
				ContentType: "image/jpeg",
			})
			if err != nil {
				m.log.WithFields(fileLogFields).WithError(err).Error("failed to upload file")
				errCh <- fmt.Errorf("failed to upload object %s: %v", objectID, err)
				cancel()
				return
			}

			// Логируем информацию о загрузке
			fileLogFields["upload_info"] = uploadInfo
			m.log.WithFields(fileLogFields).Debug("file upload details")

			// Получение URL для загруженного объекта
			url, err := m.mc.PresignedGetObject(ctx, m.config.BucketName, objectID, time.Second*24*60*60, nil)
			if err != nil {
				m.log.WithFields(fileLogFields).WithError(err).Error("failed to generate presigned URL")
				errCh <- fmt.Errorf("failed to generate URL for object %s: %v", objectID, err)
				cancel()
				return
			}

			m.log.WithFields(fileLogFields).Debug("file successfully uploaded")
			urlCh <- url.String()
		}(objectID, file)
	}

	// Ожидание завершения всех горутин и закрытие канала с URL-адресами.
	go func() {
		wg.Wait()    // Блокировка до тех пор, пока счетчик WaitGroup не станет равным 0.
		close(urlCh) // Закрытие канала с URL-адресами после завершения всех горутин.
		close(errCh)
	}()

	var urls []string
    for i := 0; i < len(data); i++ {
        select {
        case url := <-urlCh:
            urls = append(urls, url)
        case err := <-errCh:
            m.log.WithFields(logFields).WithError(err).Error("upload failed")
            return urls, err // Частичные результаты + ошибка
        }
    }

    m.log.WithFields(logFields).Info("all files successfully uploaded")
	return urls, nil
}

// GetOne получает один объект из бакета Minio по его идентификатору.
// Он принимает строку `objectID` в качестве параметра и возвращает срез байт данных объекта и ошибку, если такая возникает.
func (m *minioProvider) GetOne(ctx context.Context, objectID string) (string, error) {
	logFields := logrus.Fields{
		"object_id": objectID,
	}

	m.log.WithFields(logFields).Debug("attempting to get file URL from MinIO")

	publicURL := os.Getenv("MINIO_PUBLIC_URL")
    if publicURL == "" {
        publicURL = "/s3/"
    }

	// Получение предварительно подписанного URL для доступа к объекту Minio.
	url := fmt.Sprintf("%s%s", publicURL, objectID)


	m.log.WithFields(logFields).Debug("successfully retrieved file URL")
	return url, nil
}

// GetMany получает несколько объектов из бакета Minio по их идентификаторам.
func (m *minioProvider) GetMany(ctx context.Context, objectIDs []string) ([]string, error) {
	logFields := logrus.Fields{
		"total_objects": len(objectIDs),
	}

	m.log.WithFields(logFields).Debug("attempting to get multiple file URLs from MinIO")

	// Создание каналов для передачи URL-адресов объектов и ошибок
	urlCh := make(chan string, len(objectIDs))
	errCh := make(chan OperationError, len(objectIDs))

	var wg sync.WaitGroup                // WaitGroup для ожидания завершения всех горутин
	ctx, cancel := context.WithCancel(ctx) // Создание контекста с возможностью отмены операции
	defer cancel()                       // Отложенный вызов функции отмены контекста при завершении функции GetMany

	// Запуск горутин для получения URL-адресов каждого объекта.
	for _, objectID := range objectIDs {
		wg.Add(1)
		go func(objectID string) {
			defer wg.Done()
			fileLogFields := logrus.Fields{
				"object_id": objectID,
			}

			url, err := m.GetOne(ctx, objectID)
			if err != nil {
				m.log.WithFields(fileLogFields).WithError(err).Error("failed to get file URL")
				errCh <- OperationError{ObjectID: objectID, Error: fmt.Errorf("failed to retrieve object %s: %v", objectID, err)}
				cancel()
				return
			}
			m.log.WithFields(fileLogFields).Debug("successfully retrieved file URL")
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
    for i := 0; i < len(objectIDs); i++ {
        select {
        case url := <-urlCh:
            urls = append(urls, url)
        case err := <-errCh:
            m.log.WithFields(logFields).WithError(err.Error).Error("error occurred while getting URLs")
            return urls, err.Error // Возвращаем то, что успели собрать + ошибку
        }
    }

	m.log.WithFields(logFields).Info("successfully retrieved all file URLs")
	return urls, nil
}

// DeleteOne удаляет один объект из бакета Minio по его идентификатору.
func (m *minioProvider) DeleteOne(ctx context.Context, objectID string) error {
	logFields := logrus.Fields{
		"object_id": objectID,
	}

	m.log.WithFields(logFields).Debug("attempting to delete file from MinIO")

	// Удаление объекта из бакета Minio.
	if err := m.mc.RemoveObject(ctx, m.config.BucketName, objectID, minio.RemoveObjectOptions{}); err != nil {
		err = fmt.Errorf("failed to delete object %s: %w", objectID, err)
		m.log.WithFields(logFields).WithError(err).Error("failed to delete file")
		return err
	}

	m.log.WithFields(logFields).Info("successfully deleted file")
	return nil
}

// DeleteMany удаляет несколько объектов из бакета Minio по их идентификаторам с использованием горутин.
func (m *minioProvider) DeleteMany(ctx context.Context, objectIDs []string) error {
    logFields := logrus.Fields{
        "total_objects": len(objectIDs),
    }

    m.log.WithFields(logFields).Debug("attempting to delete multiple objects from MinIO")

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    errCh := make(chan error, 1) // Только первая ошибка
    var wg sync.WaitGroup

    for _, objectID := range objectIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            
            select {
            case <-ctx.Done():
                return
            default:
            }

            // Используем DeleteOne вместо прямого вызова RemoveObject
            if err := m.DeleteOne(ctx, id); err != nil {
                select {
                case errCh <- fmt.Errorf("failed to delete object %s: %w", id, err):
                    cancel()
                default:
                }
            }
        }(objectID)
    }

    go func() {
        wg.Wait()
        close(errCh)
    }()

    select {
    case err := <-errCh:
        m.log.WithFields(logFields).WithError(err).Error("errors occurred while deleting objects")
        return err
    default:
        m.log.WithFields(logFields).Info("all objects successfully deleted")
        return nil
    }
}

type FileData struct {
	Name string
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
