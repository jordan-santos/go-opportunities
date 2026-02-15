package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"opportunities/internal/repository"
	"opportunities/internal/schemas"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOpeningHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should return 200 when creating a valid opening", func(t *testing.T) {
		mockRepo := new(repository.OpeningRepositoryMock)
		h := New(mockRepo)

		input := schemas.Openings{
			Role:     "Go Developer",
			Company:  "Google",
			Location: "Remote",
			Link:     "https://google.com",
			Remote:   true,
			Salary:   15000,
		}
		body, _ := json.Marshal(input)

		mockRepo.On("Create", mock.AnythingOfType("*schemas.Openings")).Return(nil)

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request, _ = http.NewRequest("POST", "/opening", bytes.NewBuffer(body))
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.CreateOpeningHandler(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		mockRepo.AssertExpectations(t)
	})
}
