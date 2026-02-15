package handler

import (
	"log/slog"
	"opportunities/internal/repository"
)

type OpeningHandler struct {
	logger *slog.Logger
	repo   repository.OpeningRepository
}

func New(repo repository.OpeningRepository) *OpeningHandler {
	return &OpeningHandler{
		logger: slog.Default().With("group", "handler"),
		repo:   repo,
	}
}
