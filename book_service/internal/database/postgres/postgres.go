package postgres

import (
	"BookShop/book_service/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New(dbConfig config.DatabaseConfig) (*Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = CreateTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &Database{db: db}, nil
}
