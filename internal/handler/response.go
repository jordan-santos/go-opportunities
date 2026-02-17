package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, code int, msg string) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

func sendSuccess(c *gin.Context, op string, data interface{}) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(200, gin.H{
		"data":    data,
		"message": fmt.Sprintf("%s", op),
	})
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type openingResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	Role      string    `json:"role"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Remote    bool      `json:"remote"`
	Link      string    `json:"link"`
	Salary    int64     `json:"salary"`
}

type CreateOpeningResponse struct {
	Message string          `json:"message"`
	Data    openingResponse `json:"data"`
}

type DeleteOpeningResponse struct {
	Message string          `json:"message"`
	Data    openingResponse `json:"data"`
}
type ShowOpeningResponse struct {
	Message string          `json:"message"`
	Data    openingResponse `json:"data"`
}
type ListOpeningsResponse struct {
	Message string            `json:"message"`
	Data    []openingResponse `json:"data"`
}
type UpdateOpeningResponse struct {
	Message string          `json:"message"`
	Data    openingResponse `json:"data"`
}

type openingCSVAcceptedData struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
}

type OpeningCSVAcceptedResponse struct {
	Message string                 `json:"message"`
	Data    openingCSVAcceptedData `json:"data"`
}
