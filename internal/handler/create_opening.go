package handler

import (
	"log/slog"
	"net/http"
	"opportunities/internal/schemas"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// CreateOpeningHandler godoc
// @Summary CreateOpening
// @Schemes
// @Description Create an Opening
// @Tags Opening
// @Accept json
// @Produce json
// @Param request body CreateOpeningRequest true "Request Body"
// @Success 200 {object} CreateOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /opening [post]
func (h *OpeningHandler) CreateOpeningHandler(c *gin.Context) {
	request := CreateOpeningRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		h.logger.Error("binding request", slog.String("error", err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Error("validation ", slog.String("error", err.Error()))
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	opening := schemas.Openings{
		Role:     request.Role,
		Company:  request.Company,
		Location: request.Location,
		Remote:   *request.Remote,
		Link:     request.Link,
		Salary:   request.Salary,
	}

	if err := h.repo.Create(&opening); err != nil {
		h.logger.Error("create db ", slog.String("error", err.Error()))
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, "createOpening", opening)
}
