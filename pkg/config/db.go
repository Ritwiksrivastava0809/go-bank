package config

import (
	"database/sql"
	"fmt"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB contains the database configuration details
type DB struct {
	Username string
	Password string
	Host     string
	Port     int
	Name     string
	SslMode  string
}

// NewDB initializes a PostgreSQL database connection using the configuration
func NewDB() (*sql.DB, error) {
	// Get database config from config package
	dbConfig := DBConfig()

	// Create a connection string using the fetched config
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SslMode,
	)

	// Open the database connection
	db, err := sql.Open(constants.DBDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %w", err)
	}

	return db, nil
}
