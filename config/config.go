package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	MinioConfig  *MinioConfig
	DBConfig     *DBConfig
	ServerConfig *ServerConfig
	JWTConfig    *JWTConfig
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file. %v", err)
	}

	minioConf, err := newMinioConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := newDBConfig()
	if err != nil {
		return nil, err
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := newJWTConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		MinioConfig:  minioConf,
		DBConfig:     dbConfig,
		ServerConfig: serverConfig,
		JWTConfig:    jwtConfig,
	}, nil
}

// MinioConfig структура, обозначающая структуру .env файла
type MinioConfig struct {
	Port         string // Порт, на котором запускается сервер
	Endpoint     string // Адрес конечной точки Minio
	BucketName   string // Название конкретного бакета в Minio
	RootUser     string // Имя пользователя для доступа к Minio
	RootPassword string // Пароль для доступа к Minio
	UseSSL       bool   // Переменная, отвечающая за
}

func newMinioConfig() (*MinioConfig, error) {
	// Проверяем обязательные переменные окружения
	endpoint, endpointExists := os.LookupEnv("MINIO_ENDPOINT")
	rootUser, userExists := os.LookupEnv("MINIO_ROOT_USER")
	rootPassword, passwordExists := os.LookupEnv("MINIO_ROOT_PASSWORD")
	port, portExists := os.LookupEnv("PORT")
	bucketName, bucketExists := os.LookupEnv("MINIO_BUCKET_NAME")
	useSSL, err := getEnvAsBool("MINIO_USE_SSL")

	if err != nil {
		return nil, err
	}

	if !endpointExists || !userExists || !passwordExists || !portExists || !bucketExists {
		return nil, errors.New("incomplete MinIO configuration: missing required environment variables")
	}

	return &MinioConfig{
		Port:         port,
		Endpoint:     endpoint,
		BucketName:   bucketName,
		RootUser:     rootUser,
		RootPassword: rootPassword,
		UseSSL:       useSSL,
	}, nil
}

type DBConfig struct {
	User     string
	Password string
	DB       string
	Port     int
	Host     string
}

func newDBConfig() (*DBConfig, error) {
	user, userExists := os.LookupEnv("POSTGRES_USER")
	password, passwordExists := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbExists := os.LookupEnv("POSTGRES_DB")
	host, hostExists := os.LookupEnv("POSTGRES_HOST")
	portStr, portExists := os.LookupEnv("POSTGRES_PORT")

	if !userExists || !passwordExists || !dbExists || !hostExists || !portExists {
		return nil, errors.New("incomplete database connection information")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New("invalid POSTGRES_PORT value")
	}

	return &DBConfig{
		User:     user,
		Password: password,
		DB:       dbname,
		Port:     port,
		Host:     host,
	}, nil
}

type ServerConfig struct {
	Port        string
	AllowOrigin string
}

func newServerConfig() (*ServerConfig, error) {
	port, portExist := os.LookupEnv("SERVER_PORT")
	allowOrigin, originExist := os.LookupEnv("ALLOW_ORIGIN")
	if !portExist || !originExist {
		return nil, errors.New("incomplete server configuration: missing required environment variable")
	}

	return &ServerConfig{
		Port:        port,
		AllowOrigin: allowOrigin,
	}, nil
}

type JWTConfig struct {
	Signature string
}

func newJWTConfig() (*JWTConfig, error) {
	signature, exists := os.LookupEnv("JWT_SIGNATURE")
	if !exists {
		return nil, errors.New("jwt signature is not set")
	}

	return &JWTConfig{
		Signature: signature,
	}, nil
}

func getEnvAsBool(key string) (bool, error) {
	valueStr, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("environment variable %s is required", key)
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}

	return value, nil
}
