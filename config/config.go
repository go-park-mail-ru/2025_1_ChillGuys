package config

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

// Config объединяет все конфигурационные настройки приложения.
type Config struct {
	MinioConfig      *MinioConfig
	DBConfig         *DBConfig
	ServerConfig     *ServerConfig
	JWTConfig        *JWTConfig
	MigrationsConfig *MigrationsConfig
	CSRFConfig       *CSRFConfig
}

// NewConfig загружает переменные окружения и инициализирует все компоненты конфига.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
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

	migrationsConfig, err := newMigrationsConfig()
	if err != nil {
		return nil, err
	}

	csrfConfig, err := newCSRFConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		MinioConfig:      minioConf,
		DBConfig:         dbConfig,
		ServerConfig:     serverConfig,
		JWTConfig:        jwtConfig,
		MigrationsConfig: migrationsConfig,
		CSRFConfig:       csrfConfig,
	}, nil
}

type MinioConfig struct {
	Port         string
	Endpoint     string
	BucketName   string
	RootUser     string
	RootPassword string
	UseSSL       bool
	PublicURL    string
}

func newMinioConfig() (*MinioConfig, error) {
	endpoint, endpointExists := os.LookupEnv("MINIO_ENDPOINT")
	rootUser, userExists := os.LookupEnv("MINIO_ROOT_USER")
	rootPassword, passwordExists := os.LookupEnv("MINIO_ROOT_PASSWORD")
	port, portExists := os.LookupEnv("PORT")
	bucketName, bucketExists := os.LookupEnv("MINIO_BUCKET_NAME")
	useSSL, err := getEnvAsBool("MINIO_USE_SSL")
	publicURL, publicURLExists := os.LookupEnv("MINIO_PUBLIC_URL")

	if err != nil {
		return nil, err
	}
	if !endpointExists || !userExists || !passwordExists || !portExists || !bucketExists || !publicURLExists {
		return nil, errors.New("incomplete MinIO configuration: missing required environment variables")
	}
	return &MinioConfig{
		Port:         port,
		Endpoint:     endpoint,
		BucketName:   bucketName,
		RootUser:     rootUser,
		RootPassword: rootPassword,
		UseSSL:       useSSL,
		PublicURL:    publicURL,
	}, nil
}

type DBConfig struct {
	User            string
	Password        string
	DB              string
	Port            int
	Host            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
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

	maxOpenConns := 25
	if val, exists := os.LookupEnv("DB_MAX_OPEN_CONNS"); exists {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxOpenConns = parsed
		}
	}
	maxIdleConns := 25
	if val, exists := os.LookupEnv("DB_MAX_IDLE_CONNS"); exists {
		if parsed, err := strconv.Atoi(val); err == nil {
			maxIdleConns = parsed
		}
	}
	connMaxLifetime := 5 * time.Minute
	if val, exists := os.LookupEnv("DB_CONN_MAX_LIFETIME"); exists {
		if parsed, err := strconv.Atoi(val); err == nil {
			connMaxLifetime = time.Duration(parsed) * time.Minute
		}
	}

	return &DBConfig{
		User:            user,
		Password:        password,
		DB:              dbname,
		Port:            port,
		Host:            host,
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}, nil
}

func ConfigureDB(db *sql.DB, cfg *DBConfig) {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
}

type ServerConfig struct {
	Port               string
	AllowOrigin        string
	AllowMethods       string
	AllowHeaders       string
	AllowCredentials   string
	WriteTimeout       time.Duration
	ReadTimeout        time.Duration
	IdleTimeout        time.Duration
	MaxMultipartMemory int64
	AvatarKey          string
}

func newServerConfig() (*ServerConfig, error) {
	port, portExist := os.LookupEnv("SERVER_PORT")
	allowOrigin, originExist := os.LookupEnv("ALLOW_ORIGIN")
	if !portExist || !originExist {
		return nil, errors.New("incomplete server configuration: missing required environment variable")
	}

	allowMethods := getEnvWithDefault("ALLOW_METHODS", "POST,GET,PUT,DELETE,OPTIONS")
	allowHeaders := getEnvWithDefault("ALLOW_HEADERS", "Content-Type,X-CSRF-Token")
	allowCredentials := getEnvWithDefault("ALLOW_CREDENTIALS", "true")

	writeTimeout := getEnvAsDuration("SERVER_WRITE_TIMEOUT", 10*time.Second)
	readTimeout := getEnvAsDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	idleTimeout := getEnvAsDuration("SERVER_IDLE_TIMEOUT", 30*time.Second)

	maxMultipartMemory := int64(10 << 20) // 10 MB по умолчанию
	if val, exists := os.LookupEnv("SERVER_MAX_MULTIPART_MEMORY"); exists {
		if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
			maxMultipartMemory = parsed
		}
	}

	avatarKey := "file"
	if val, exists := os.LookupEnv("SERVER_AVATAR_KEY"); exists {
		avatarKey = val
	}

	return &ServerConfig{
		Port:               port,
		AllowOrigin:        allowOrigin,
		AllowMethods:       allowMethods,
		AllowHeaders:       allowHeaders,
		AllowCredentials:   allowCredentials,
		WriteTimeout:       writeTimeout,
		ReadTimeout:        readTimeout,
		IdleTimeout:        idleTimeout,
		MaxMultipartMemory: maxMultipartMemory,
		AvatarKey:          avatarKey,
	}, nil
}

type JWTConfig struct {
	Signature     string
	TokenLifeSpan time.Duration
}

func newJWTConfig() (*JWTConfig, error) {
	signature, exists := os.LookupEnv("JWT_SIGNATURE")
	if !exists {
		return nil, errors.New("jwt signature is not set")
	}

	tokenLifeSpan := getEnvAsDuration("JWT_TOKEN_LIFESPAN", 24*time.Hour)

	return &JWTConfig{
		Signature:     signature,
		TokenLifeSpan: tokenLifeSpan,
	}, nil
}

type MigrationsConfig struct {
	Path string
}

func newMigrationsConfig() (*MigrationsConfig, error) {
	path, exists := os.LookupEnv("MIGRATIONS_PATH")
	if !exists {
		return nil, errors.New("MIGRATIONS_PATH is not set")
	}
	return &MigrationsConfig{
		Path: path,
	}, nil
}

type CSRFConfig struct {
	SecretKey    string
	TokenExpiry  time.Duration
	CookieName   string
	SecureCookie bool
}

func newCSRFConfig() (*CSRFConfig, error) {
	secretKey, exists := os.LookupEnv("CSRF_SECRET_KEY")
	if !exists {
		return nil, errors.New("CSRF_SECRET_KEY is not set")
	}

	tokenExpiry := getEnvAsDuration("CSRF_TOKEN_EXPIRY", 24*time.Hour)
	cookieName := getEnvWithDefault("CSRF_COOKIE_NAME", "_csrf")
	secureCookie := getEnvWithDefault("CSRF_SECURE_COOKIE", "true") == "true"

	return &CSRFConfig{
		SecretKey:    secretKey,
		TokenExpiry:  tokenExpiry,
		CookieName:   cookieName,
		SecureCookie: secureCookie,
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

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists {
		if parsed, err := time.ParseDuration(val); err == nil {
			return parsed
		}
	}
	return defaultVal
}

func getEnvWithDefault(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
