package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"opportunities/internal/auth"
	"opportunities/internal/middleware"
	"opportunities/internal/repository"
	"opportunities/internal/schemas"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOpeningHandler_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockRepo := new(repository.OpeningRepositoryMock)
	h := New(mockRepo)

	r := gin.Default()
	r.Use(middleware.Auth())
	r.POST("/opening", h.CreateOpeningHandler)

	input := schemas.Openings{
		Role:     "Go Developer",
		Company:  "Google",
		Location: "USA",
		Remote:   true,
		Link:     "https://google.com",
		Salary:   15000,
	}
	body, _ := json.Marshal(input)

	t.Run("Should return 401 Unauthorized when token is missing", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/opening", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()

		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)

		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Should return 200 Created when token is valid", func(t *testing.T) {
		token, _ := auth.GenerateToken("test@test.com")

		req, _ := http.NewRequest("POST", "/opening", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		recorder := httptest.NewRecorder()

		mockRepo.On("Create", mock.AnythingOfType("*schemas.Openings")).Return(nil).Once()

		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockRepo.AssertExpectations(t)
	})
}
