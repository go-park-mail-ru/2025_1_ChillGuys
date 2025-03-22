package utils

import (
	"errors"
	"fmt"
	"os"
)

func GetConnectionString() (string, error) {
	user, userExists := os.LookupEnv("POSTGRES_USER")
	password, passwordExists := os.LookupEnv("POSTGRES_PASSWORD")
	dbname, dbExists := os.LookupEnv("POSTGRES_DB")
	host, hostExists := os.LookupEnv("POSTGRES_HOST")
	port, portExists := os.LookupEnv("POSTGRES_PORT")

	if !userExists || !passwordExists || !dbExists || !hostExists || !portExists {
		return "", errors.New("incomplete database connection information")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	return connStr, nil
}
