package config

import (
	"os"
	"strings"
)

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	ClientID string
}

func LoadKafkaConfig() KafkaConfig {
	brokersRaw := strings.TrimSpace(os.Getenv("KAFKA_BROKERS"))
	if brokersRaw == "" {
		brokersRaw = "localhost:9092"
	}

	topic := strings.TrimSpace(os.Getenv("KAFKA_TOPIC_FEEDBACK"))
	if topic == "" {
		topic = "feedback-opening-v1"
	}

	clientID := strings.TrimSpace(os.Getenv("KAFKA_CLIENT_ID"))
	if clientID == "" {
		clientID = "opportunities-api"
	}

	brokers := strings.Split(brokersRaw, ",")
	for i := range brokers {
		brokers[i] = strings.TrimSpace(brokers[i])
	}

	return KafkaConfig{
		Brokers:  brokers,
		Topic:    topic,
		ClientID: clientID,
	}
}
