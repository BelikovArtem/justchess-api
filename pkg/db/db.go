package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "embed"

	_ "github.com/lib/pq"
)

//go:embed schema.sql
var schema string

var DB *sql.DB

// OpenDB opens a new connection with the database,
// using env variables to format a connection string
func OpenDB() error {
	connectStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}
	slog.Info("Database connected successfully.")
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	return nil
}

// CreateTables reads the embeded schema file and executes all queries from this file.
func CreateTables() error {
	if DB == nil {
		return fmt.Errorf("database connection is not opened")
	}
	_, err := DB.Query(string(schema))
	if err != nil {
		return err
	}
	slog.Info("Tables created successfully.")
	return nil
}

func CloseDB() error {
	err := DB.Close()
	return err
}
