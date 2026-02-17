package handler

import (
	"io"
	"log/slog"
	"net/http"

	csvutil "opportunities/internal/csv"
	"opportunities/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @BasePath /api/v1

// CreateOpeningCSVHandler godoc
// @Summary Create openings from CSV
// @Description Upload a CSV and process openings asynchronously
// @Tags Opening
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV file"
// @Success 202 {object} OpeningCSVAcceptedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 503 {object} ErrorResponse
// @Security BearerAuth
// @Router /opening/csv [post]
func (h *OpeningHandler) CreateOpeningCSVHandler(c *gin.Context) {
	if h.csvService == nil {
		sendError(c, http.StatusServiceUnavailable, "csv service unavailable")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		sendError(c, http.StatusBadRequest, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.logger.Error("CreateOpeningCSVHandler open file", slog.String("error", err.Error()))
		sendError(c, http.StatusBadRequest, "invalid file")
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		h.logger.Error("CreateOpeningCSVHandler read file", slog.String("error", err.Error()))
		sendError(c, http.StatusBadRequest, "invalid file")
		return
	}

	if err := csvutil.ValidateHeader(content); err != nil {
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	requestID := uuid.NewString()

	err = h.csvService.Enqueue(service.OpeningCSVJob{
		RequestID: requestID,
		Content:   content,
	})
	if err != nil {
		if err == service.ErrCSVQueueFull {
			sendError(c, http.StatusServiceUnavailable, err.Error())
			return
		}

		h.logger.Error("CreateOpeningCSVHandler enqueue job", slog.String("error", err.Error()))
		sendError(c, http.StatusInternalServerError, "failed to enqueue csv processing")
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusAccepted, gin.H{
		"message": "openingCsvAccepted",
		"data": gin.H{
			"request_id": requestID,
			"status":     "accepted",
		},
	})
}
