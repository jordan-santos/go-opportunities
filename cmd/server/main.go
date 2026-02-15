package main

import (
	"log/slog"
	"opportunities/config"
	_ "opportunities/docs"
	"opportunities/internal/router"
)

// @title Opportunities API
// @version 1.0
// @description API para gerenciamento de vagas de emprego.
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := config.Init()

	if err != nil {
		slog.Error("Error initializing config", slog.String("error", err.Error()))
		return
	}

	db := config.GetSQLite()

	router.Initialize(db)
}
