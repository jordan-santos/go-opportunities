package handler

import (
	"net/http"
	"opportunities/schemas"

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
// @Router /opening [post]
func CreateOpeningHandler(c *gin.Context) {
	request := CreateOpeningRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		logger.Errorf("binding request error: %s", err)
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("validation error: %v", err)
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

	if err := db.Create(&opening).Error; err != nil {
		logger.Errorf("create db error: %s", err)
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, "createOpening", opening)
}
