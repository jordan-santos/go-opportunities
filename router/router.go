package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Initialize(db *gorm.DB) {
	router := gin.Default()

	initializeRoutes(router, db)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
