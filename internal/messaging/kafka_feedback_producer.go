package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducerConfig struct {
	Brokers  []string
	Topic    string
	ClientID string
}

type KafkaFeedbackProducer struct {
	writer  *kafka.Writer
	brokers []string
	topic   string
}

func NewKafkaFeedbackProducer(cfg KafkaProducerConfig) *KafkaFeedbackProducer {
	producer := &KafkaFeedbackProducer{
		brokers: cfg.Brokers,
		topic:   cfg.Topic,
	}

	producer.writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		AllowAutoTopicCreation: true,
		Async:        false,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Transport: &kafka.Transport{
			ClientID: cfg.ClientID,
		},
	}

	// Best-effort startup topic ensure to reduce first-message race conditions.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := producer.ensureTopic(ctx, 1, 1); err != nil {
		for i := 0; i < 4; i++ {
			time.Sleep(1 * time.Second)
			if retryErr := producer.ensureTopic(ctx, 1, 1); retryErr == nil {
				break
			} else if i == 3 {
				fmt.Println("kafka topic ensure warning:", retryErr.Error())
			}
		}
	}

	return producer
}

func (p *KafkaFeedbackProducer) PublishOpeningCSVFeedback(ctx context.Context, feedback OpeningCSVFeedback) error {
	payload, err := json.Marshal(feedback)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(feedback.RequestID),
		Value: payload,
		Time:  time.Now().UTC(),
	}

	var lastErr error
	for i := 0; i < 4; i++ {
		err = p.writer.WriteMessages(ctx, msg)
		if err == nil {
			return nil
		}

		lastErr = err
		if !isUnknownTopicOrPartitionErr(err) {
			return err
		}

		_ = p.ensureTopic(ctx, 1, 1)
		time.Sleep(500 * time.Millisecond)
	}

	return lastErr
}

func (p *KafkaFeedbackProducer) ensureTopic(ctx context.Context, partitions int, replicationFactor int) error {
	if len(p.brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}

	conn, err := kafka.DialContext(ctx, "tcp", p.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	parts, err := conn.ReadPartitions(p.topic)
	if err == nil && len(parts) > 0 {
		return nil
	}

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	ctrlAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	ctrlConn, err := kafka.DialContext(ctx, "tcp", ctrlAddr)
	if err != nil {
		return err
	}
	defer ctrlConn.Close()

	err = ctrlConn.CreateTopics(kafka.TopicConfig{
		Topic:             p.topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "already exists") {
		return err
	}

	return nil
}

func isUnknownTopicOrPartitionErr(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unknown topic or partition")
}
