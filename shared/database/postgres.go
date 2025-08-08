package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"ua/shared/logger"
)

type DB struct {
	*sql.DB
}

func NewPostgresDB(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to PostgreSQL database")
	return &DB{db}, nil
}

func (db *DB) Close() error {
	logger.Info("Closing PostgreSQL connection")
	return db.DB.Close()
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	logger.Debug("Executing query", zap.String("query", query))
	return db.DB.Exec(query, args...)
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	logger.Debug("Executing query", zap.String("query", query))
	return db.DB.Query(query, args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	logger.Debug("Executing query", zap.String("query", query))
	return db.DB.QueryRow(query, args...)
}
