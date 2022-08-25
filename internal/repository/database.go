package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

const (
	driver    = "postgres"
	dsnFormat = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
)

var (
	errEmptyDatabaseUsername = errors.New("empty database username")
	errEmptyDatabasePassword = errors.New("empty database password")
	errEmptyDatabaseHostname = errors.New("empty database hostname")
	errInvalidDatabasePort   = errors.New("invalid database port")
	errEmptyDatabaseName     = errors.New("empty database name")
)

// NewDB creates a new connection to the database.
func NewDB() (*sql.DB, error) {
	if os.Getenv("DB_USER") == "" {
		return nil, errEmptyDatabaseUsername
	}

	if os.Getenv("DB_PASSWORD") == "" {
		return nil, errEmptyDatabasePassword
	}

	if os.Getenv("DB_HOST") == "" {
		return nil, errEmptyDatabaseHostname
	}

	_, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, errInvalidDatabasePort
	}

	if os.Getenv("DB_NAME") == "" {
		return nil, errEmptyDatabaseName
	}

	dsn := fmt.Sprintf(
		dsnFormat,
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	return sql.Open(driver, dsn)
}
