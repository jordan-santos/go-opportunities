package messaging

import (
	"context"
	"time"
)

type OpeningCSVFeedback struct {
	RequestID      string    `json:"request_id"`
	Status         string    `json:"status"`
	TotalRows      int       `json:"total_rows"`
	ProcessedRows  int       `json:"processed_rows"`
	DurationMS     int64     `json:"duration_ms"`
	ErrorCount     int       `json:"error_count"`
	FirstErrorLine int       `json:"first_error_line"`
	Message        string    `json:"message"`
	Timestamp      time.Time `json:"timestamp"`
}

type FeedbackProducer interface {
	PublishOpeningCSVFeedback(ctx context.Context, feedback OpeningCSVFeedback) error
}
