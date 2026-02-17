package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"opportunities/internal/auth"
	"opportunities/internal/middleware"
	"opportunities/internal/repository"
	"opportunities/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOpeningCSVHandler_WithAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should return 401 Unauthorized when token is missing", func(t *testing.T) {
		mockRepo := new(repository.OpeningRepositoryMock)
		csvService := service.NewOpeningCSVService(mockRepo, nil, 1)
		h := New(mockRepo, csvService)
		r := gin.Default()
		r.Use(middleware.Auth())
		r.POST("/opening/csv", h.CreateOpeningCSVHandler)

		body, contentType := newCSVMultipartBody(t, "file", "openings.csv", "role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\n")
		req, _ := http.NewRequest("POST", "/opening/csv", body)
		req.Header.Set("Content-Type", contentType)

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	t.Run("Should return 202 Accepted when token and csv are valid", func(t *testing.T) {
		mockRepo := new(repository.OpeningRepositoryMock)
		csvService := service.NewOpeningCSVService(mockRepo, nil, 1)
		h := New(mockRepo, csvService)
		r := gin.Default()
		r.Use(middleware.Auth())
		r.POST("/opening/csv", h.CreateOpeningCSVHandler)

		token, _ := auth.GenerateToken("test@test.com")
		body, contentType := newCSVMultipartBody(t, "file", "openings.csv", "role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\n")
		req, _ := http.NewRequest("POST", "/opening/csv", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+token)

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusAccepted, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "openingCsvAccepted")
		assert.Contains(t, recorder.Body.String(), "request_id")
	})

	t.Run("Should return 400 for invalid csv header", func(t *testing.T) {
		mockRepo := new(repository.OpeningRepositoryMock)
		csvService := service.NewOpeningCSVService(mockRepo, nil, 1)
		h := New(mockRepo, csvService)
		r := gin.Default()
		r.Use(middleware.Auth())
		r.POST("/opening/csv", h.CreateOpeningCSVHandler)

		token, _ := auth.GenerateToken("test@test.com")
		body, contentType := newCSVMultipartBody(t, "file", "openings.csv", "role,company,location,link,salary\nGo Dev,Acme,BR,https://acme.com,1000\n")
		req, _ := http.NewRequest("POST", "/opening/csv", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+token)

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("Should return 503 when queue is full", func(t *testing.T) {
		mockRepo := new(repository.OpeningRepositoryMock)
		csvService := service.NewOpeningCSVService(mockRepo, nil, 0)
		h := New(mockRepo, csvService)
		r := gin.Default()
		r.Use(middleware.Auth())
		r.POST("/opening/csv", h.CreateOpeningCSVHandler)

		token, _ := auth.GenerateToken("test@test.com")
		body, contentType := newCSVMultipartBody(t, "file", "openings.csv", "role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\n")
		req, _ := http.NewRequest("POST", "/opening/csv", body)
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", "Bearer "+token)

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusServiceUnavailable, recorder.Code)
	})
}

func newCSVMultipartBody(t *testing.T, fieldName, fileName, fileContent string) (*bytes.Buffer, string) {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("CreateFormFile error: %v", err)
	}

	_, err = part.Write([]byte(fileContent))
	if err != nil {
		t.Fatalf("write multipart file error: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("writer close error: %v", err)
	}

	return body, writer.FormDataContentType()
}
