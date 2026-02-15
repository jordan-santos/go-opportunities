package config

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() error {
	var err error

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	db, err = InitializeSQLite()
	if err != nil {
		return fmt.Errorf("error initializing sqlite database: %v", err)
	}

	return nil
}

func GetSQLite() *gorm.DB {
	return db
}

func GetLogger(p string) *slog.Logger {
	return slog.Default().With("source", p)
}
