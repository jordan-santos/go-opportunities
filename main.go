package main

import (
	"opportunities/config"
	"opportunities/router"
)

func main() {
	err := config.Init()
	logger := config.GetLogger("main")

	if err != nil {
		logger.Errorf("Error initializing config: %v", err)
		return
	}

	db := config.GetSQLite()

	router.Initialize(db)
}
