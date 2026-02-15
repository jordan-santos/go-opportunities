package main

import (
	"log/slog"
	"opportunities/config"
	"opportunities/internal/router"
)

func main() {
	err := config.Init()

	if err != nil {
		slog.Error("Error initializing config", slog.String("error", err.Error()))
		return
	}

	db := config.GetSQLite()

	router.Initialize(db)
}
