package config

import (
	"fmt"
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
		logger.Info("Creating new sqlite database")
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
		logger.Errorf("Error opening SQLite database: %v", err)
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	err = db.AutoMigrate(&schemas.Openings{})
	if err != nil {
		logger.Errorf("Error auto-migrating SQLite database: %v", err)
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	return db, nil
}
