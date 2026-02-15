package router

import (
	"opportunities/docs"
	"opportunities/internal/handler"
	"opportunities/internal/middleware"
	"opportunities/internal/repository"

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

	router.POST(basePath+"/login", h.LoginHandler)

	v1Public := router.Group(basePath)
	{
		v1Public.GET("/opening", h.ShowOpeningHandler)
		v1Public.GET("/openings", h.ListOpeningHandler)
	}

	v1Protected := router.Group(basePath)
	v1Protected.Use(middleware.Auth())
	{
		v1Protected.POST("/opening", h.CreateOpeningHandler)
		v1Protected.PUT("/opening", h.UpdateOpeningHandler)
		v1Protected.DELETE("/opening", h.DeleteOpeningHandler)
	}

	// swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
