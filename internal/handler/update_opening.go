package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// UpdateOpeningHandler godoc
// @Summary Update opening
// @Description Update a job opening
// @Tags Opening
// @Accept json
// @Produce json
// @Param id query string true "Opening Identification"
// @Param opening body UpdateOpeningRequest true "Opening data to Update"
// @Success 200 {object} UpdateOpeningResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /opening [put]
func (h *OpeningHandler) UpdateOpeningHandler(c *gin.Context) {
	request := UpdateOpeningRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		h.logger.Error("UpdateOpeningHandler parse request body", slog.String("error", err.Error()))
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Error("UpdateOpeningHandler validate request", slog.String("error", err.Error()))
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Query("id")
	if id == "" {
		sendError(c, http.StatusBadRequest, errParamIsRequired("id", "queryParam").Error())
		return
	}

	opening, err := h.repo.Get(id)
	if err != nil {
		sendError(c, http.StatusNotFound, fmt.Sprintf("opening %s not found", id))
		return
	}

	if request.Role != "" {
		opening.Role = request.Role
	}

	if request.Company != "" {
		opening.Company = request.Company
	}

	if request.Location != "" {
		opening.Location = request.Location
	}

	if request.Remote != nil {
		opening.Remote = *request.Remote
	}

	if request.Link != "" {
		opening.Link = request.Link
	}

	if request.Salary > 0 {
		opening.Salary = request.Salary
	}

	if err := h.repo.Update(&opening); err != nil {
		h.logger.Error("UpdateOpeningHandler save opening", slog.String("error", err.Error()))
		sendError(c, http.StatusInternalServerError, err.Error())
	}

	sendSuccess(c, "updateOpening", opening)
}
