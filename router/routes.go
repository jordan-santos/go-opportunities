package router

import (
	"opportunities/handler"

	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/opening", func(c *gin.Context) {
			handler.ShowOpeningHandler(c)
		})
		v1.POST("/opening", func(c *gin.Context) {
			handler.CreateOpeningHandler(c)
		})
		v1.DELETE("/opening", func(c *gin.Context) {
			handler.DeleteOpeningHandler(c)
		})
		v1.PUT("/opening", func(c *gin.Context) {
			handler.UpdateOpeningHandler(c)
		})
		v1.GET("/openings", func(c *gin.Context) {
			handler.ListOpeningHandler(c)
		})
	}
}
