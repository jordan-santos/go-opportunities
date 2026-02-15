package router

import (
	"opportunities/docs"
	"opportunities/handler"
	"opportunities/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func initializeRoutes(router *gin.Engine, db *gorm.DB) {
	repo := repository.New(db)
	h := handler.New(repo)

	basePath := "/api/v1"
	docs.SwaggerInfo.BasePath = basePath

	v1 := router.Group(basePath)
	{
		v1.GET("/opening", h.ShowOpeningHandler)
		v1.POST("/opening", h.CreateOpeningHandler)
		v1.PUT("/opening", h.UpdateOpeningHandler)
		v1.DELETE("/opening", h.DeleteOpeningHandler)
		v1.GET("/openings", h.ListOpeningHandler)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
