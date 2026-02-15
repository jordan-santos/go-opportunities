package handler

import (
	"opportunities/config"
	"opportunities/repository"
)

type OpeningHandler struct {
	logger *config.Logger
	repo   repository.OpeningRepository
}

func New(repo repository.OpeningRepository) *OpeningHandler {
	return &OpeningHandler{
		logger: config.GetLogger("handler"),
		repo:   repo,
	}
}
