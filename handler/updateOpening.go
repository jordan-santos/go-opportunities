package handler

import (
	"net/http"
	"opportunities/schemas"

	"github.com/gin-gonic/gin"
)

func UpdateOpeningHandler(c *gin.Context) {
	request := UpdateOpeningRequest{}

	err := c.BindJSON(&request)
	if err != nil {
		logger.Errorf("UpdateOpeningHandler parse request body failed: %v", err)
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		logger.Errorf("UpdateOpeningHandler validate request failed: %v", err)
		sendError(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Query("id")
	if id == "" {
		sendError(c, http.StatusBadRequest, errParamIsRequired("id", "queryParam").Error())
		return
	}

	opening := schemas.Openings{}
	if err := db.First(&opening, id).Error; err != nil {
		sendError(c, http.StatusBadRequest, "opening not found")
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

	if err := db.Save(&opening).Error; err != nil {
		logger.Errorf("UpdateOpeningHandler save opening failed: %v", err)
		sendError(c, http.StatusInternalServerError, err.Error())
	}

	sendSuccess(c, "updateOpening", opening)
}
