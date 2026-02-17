package main

import (
	"context"
	"log/slog"
	"opportunities/config"
	_ "opportunities/docs"
	"opportunities/internal/messaging"
	"opportunities/internal/repository"
	"opportunities/internal/router"
	"opportunities/internal/service"
)

// @title Opportunities API
// @version 1.0
// @description API para gerenciamento de vagas de emprego.
// @host localhost:8080
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := config.Init()

	if err != nil {
		slog.Error("Error initializing config", slog.String("error", err.Error()))
		return
	}

	db := config.GetSQLite()
	repo := repository.New(db)
	kafkaConfig := config.LoadKafkaConfig()

	feedbackProducer := messaging.NewKafkaFeedbackProducer(messaging.KafkaProducerConfig{
		Brokers:  kafkaConfig.Brokers,
		Topic:    kafkaConfig.Topic,
		ClientID: kafkaConfig.ClientID,
	})

	csvService := service.NewOpeningCSVService(repo, feedbackProducer, 100)
	csvService.Start(context.Background())

	router.Initialize(db, csvService)
}
