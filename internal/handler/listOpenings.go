package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// ListOpeningHandler godoc
// @Summary List openings
// @Description List all job openings
// @Tags Openings
// @Accept json
// @Produce json
// @Success 200 {object} ListOpeningsResponse
// @Failure 500 {object} ErrorResponse
// @Router /openings [get]
func (h *OpeningHandler) ListOpeningHandler(c *gin.Context) {
	openings, err := h.repo.List()
	if err != nil {
		sendError(c, http.StatusInternalServerError, "error getting openings")
		return
	}

	sendSuccess(c, "openings", openings)
}
