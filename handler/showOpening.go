package handler

import (
	"fmt"
	"net/http"
	"opportunities/schemas"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// ShowOpeningHandler godoc
// @Summary Show opening
// @Description Show a job opening
// @Tags Opening
// @Accept json
// @Produce json
// @Param id query string true "Opening identification"
// @Success 200 {object} ShowOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /opening [get]
func ShowOpeningHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		sendError(c, http.StatusBadRequest,
			errParamIsRequired("id", "queryParameter").Error())
		return
	}

	opening := schemas.Openings{}
	if err := db.First(&opening, id).Error; err != nil {
		sendError(c, http.StatusNotFound, fmt.Sprintf("opening %s not found", id))
		return
	}

	sendSuccess(c, "opening", opening)
}
