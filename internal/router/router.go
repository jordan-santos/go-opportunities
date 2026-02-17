package router

import (
	"opportunities/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB, csvService *service.OpeningCSVService) {
	router := gin.Default()

	initializeRoutes(router, db, csvService)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
