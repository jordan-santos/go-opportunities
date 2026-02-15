package handler

import (
	"net/http"
	"opportunities/schemas"

	"github.com/gin-gonic/gin"
)

func ListOpeningHandler(c *gin.Context) {
	openings := []schemas.Openings{}

	if err := db.Find(&openings).Error; err != nil {
		sendError(c, http.StatusInternalServerError, "error getting openings")
		return
	}

	sendSuccess(c, "openings", openings)
}
