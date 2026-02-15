package handler

import (
	"fmt"
	"net/http"

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
func (h *OpeningHandler) DeleteOpeningHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		sendError(c, http.StatusBadRequest,
			errParamIsRequired("id", "queryParameter").Error())
		return
	}

	if err := h.repo.Delete(id).Error; err != nil {
		sendError(c, http.StatusInternalServerError,
			fmt.Sprintf("error deleting opening %s", id))
		return
	}

	sendSuccess(c, "deleteOpening", id)
}
