package config

import (
	"fmt"
	"log/slog"
	"opportunities/schemas"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitializeSQLite() (*gorm.DB, error) {
	logger := GetLogger("sqlite")
	dbPath := "./db/openings.db"

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		logger.Info("database file not found, creating a new one", slog.String("path", dbPath))

		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}

		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create database file: %w", err)
		}

		err = file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close database file: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Error("failed to open sqlite database", slog.Any("error", err))
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	err = db.AutoMigrate(&schemas.Openings{})
	if err != nil {
		logger.Error("sqlite auto-migration failed", slog.Any("error", err))
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	logger.Info("sqlite database initialized and migrated successfully")
	return db, nil
}
