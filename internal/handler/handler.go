package handler

import (
	"log/slog"
	"opportunities/internal/repository"
	"opportunities/internal/service"
)

type OpeningHandler struct {
	logger     *slog.Logger
	repo       repository.OpeningRepository
	csvService *service.OpeningCSVService
}

func New(repo repository.OpeningRepository, csvService *service.OpeningCSVService) *OpeningHandler {
	return &OpeningHandler{
		logger:     slog.Default().With("group", "handler"),
		repo:       repo,
		csvService: csvService,
	}
}
