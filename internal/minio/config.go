package minio

import (
	"os"
	"strconv"
)

// Config структура, обозначающая структуру .env файла
type Config struct {
	Port              string // Порт, на котором запускается сервер
	MinioEndpoint     string // Адрес конечной точки Minio
	BucketName        string // Название конкретного бакета в Minio
	MinioRootUser     string // Имя пользователя для доступа к Minio
	MinioRootPassword string // Пароль для доступа к Minio
	MinioUseSSL       bool   // Переменная, отвечающая за
}

func NewConfig() *Config {
	// Устанавливаем конфигурационные параметры
	return &Config{
		Port:              getEnv("PORT", "8080"),
		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
		BucketName:        getEnv("MINIO_BUCKET_NAME", "defaultBucket"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", "root"),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", "minio_password"),
		MinioUseSSL:       getEnvAsBool("MINIO_USE_SSL", false),
	}
}

// getEnv считывает значение переменной окружения или возвращает значение по умолчанию, если переменная не установлена
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsBool считывает значение переменной окружения как булево или возвращает значение по умолчанию, если переменная не установлена или не может быть преобразована в булево
func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
