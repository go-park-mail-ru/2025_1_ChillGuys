package utils

import (
	"fmt"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
)

func GetConnectionString(conf *config.DBConfig) (string, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DB,
	)

	return connStr, nil
}
