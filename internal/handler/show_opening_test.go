package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"opportunities/internal/repository"
	"opportunities/internal/schemas"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShowOpeningHandler_Table(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		idQuery      string
		mockBehavior func(m *repository.OpeningRepositoryMock)
		expectedCode int
	}{
		{
			name:    "Success - Opening Found",
			idQuery: "1",
			mockBehavior: func(m *repository.OpeningRepositoryMock) {
				m.On("Get", "1").Return(schemas.Openings{Role: "Go Developer"}, nil).Once()
			},
			expectedCode: http.StatusOK,
		},
		{
			name:    "Error - ID not provided",
			idQuery: "",
			mockBehavior: func(m *repository.OpeningRepositoryMock) {
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "Error - Opening Not Found",
			idQuery: "999",
			mockBehavior: func(m *repository.OpeningRepositoryMock) {
				m.On("Get", "999").Return(schemas.Openings{}, errors.New("not found")).Once()
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(repository.OpeningRepositoryMock)
			tt.mockBehavior(mockRepo)
			h := New(mockRepo, nil)

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			req, _ := http.NewRequest("GET", "/opening?id="+tt.idQuery, nil)
			ctx.Request = req

			h.ShowOpeningHandler(ctx)

			assert.Equal(t, tt.expectedCode, recorder.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
