package handler

import (
	"fmt"
	"net/http"
	"opportunities/schemas"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// DeleteOpeningHandler godoc
// @Summary DeleteOpening
// @Schemes
// @Description Delete an Opening
// @Tags Opening
// @Accept json
// @Produce json
// @Param id query string true "Opening identification"
// @Success 200 {object} DeleteOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opening [delete]
func DeleteOpeningHandler(c *gin.Context) {
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

	if err := db.Delete(&opening).Error; err != nil {
		sendError(c, http.StatusInternalServerError,
			fmt.Sprintf("error deleting opening %s", id))
		return
	}

	sendSuccess(c, "deleteOpening", opening)
}
