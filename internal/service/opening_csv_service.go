package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	csvutil "opportunities/internal/csv"
	"opportunities/internal/messaging"
	"opportunities/internal/repository"
)

var ErrCSVQueueFull = errors.New("csv processing queue is full")

type OpeningCSVJob struct {
	RequestID string
	Content   []byte
}

type OpeningCSVService struct {
	logger   *slog.Logger
	repo     repository.OpeningRepository
	producer messaging.FeedbackProducer
	jobs     chan OpeningCSVJob
}

func NewOpeningCSVService(repo repository.OpeningRepository, producer messaging.FeedbackProducer, queueSize int) *OpeningCSVService {
	return &OpeningCSVService{
		logger:   slog.Default().With("group", "opening_csv_service"),
		repo:     repo,
		producer: producer,
		jobs:     make(chan OpeningCSVJob, queueSize),
	}
}

func (s *OpeningCSVService) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("opening csv service stopped")
				return
			case job := <-s.jobs:
				s.processJob(ctx, job)
			}
		}
	}()
}

func (s *OpeningCSVService) Enqueue(job OpeningCSVJob) error {
	select {
	case s.jobs <- job:
		return nil
	default:
		return ErrCSVQueueFull
	}
}

func (s *OpeningCSVService) processJob(ctx context.Context, job OpeningCSVJob) {
	startTime := time.Now()
	logger := s.logger.With("request_id", job.RequestID)

	logger.Info("starting csv processing")

	parsedRows, rowErrors, err := csvutil.ParseAndValidate(job.Content)
	if err != nil {
		logger.Error("failed to parse csv", slog.String("error", err.Error()))
		s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
			RequestID:      job.RequestID,
			Status:         "error",
			TotalRows:      0,
			ProcessedRows:  0,
			DurationMS:     time.Since(startTime).Milliseconds(),
			ErrorCount:     1,
			FirstErrorLine: 0,
			Message:        err.Error(),
			Timestamp:      time.Now().UTC(),
		})
		return
	}

	totalRows := len(parsedRows) + len(rowErrors)
	if len(rowErrors) > 0 {
		for _, rowErr := range rowErrors {
			logger.Error("csv validation error",
				slog.Int("line_number", rowErr.LineNumber),
				slog.String("error", rowErr.Message))
		}

		s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
			RequestID:      job.RequestID,
			Status:         "error",
			TotalRows:      totalRows,
			ProcessedRows:  0,
			DurationMS:     time.Since(startTime).Milliseconds(),
			ErrorCount:     len(rowErrors),
			FirstErrorLine: rowErrors[0].LineNumber,
			Message:        "csv validation failed",
			Timestamp:      time.Now().UTC(),
		})
		return
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		logger.Error("failed to begin transaction", slog.String("error", err.Error()))
		s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
			RequestID:      job.RequestID,
			Status:         "error",
			TotalRows:      totalRows,
			ProcessedRows:  0,
			DurationMS:     time.Since(startTime).Milliseconds(),
			ErrorCount:     1,
			FirstErrorLine: 0,
			Message:        "failed to begin transaction",
			Timestamp:      time.Now().UTC(),
		})
		return
	}

	processed := 0
	for _, row := range parsedRows {
		opening := row.Opening
		if err := s.repo.CreateWithTx(tx, &opening); err != nil {
			logger.Error("failed to insert csv row",
				slog.Int("line_number", row.LineNumber),
				slog.String("error", err.Error()))
			tx.Rollback()
			logger.Error("transaction rolled back", slog.Bool("transaction_rolled_back", true))

			s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
				RequestID:      job.RequestID,
				Status:         "error",
				TotalRows:      totalRows,
				ProcessedRows:  0,
				DurationMS:     time.Since(startTime).Milliseconds(),
				ErrorCount:     1,
				FirstErrorLine: row.LineNumber,
				Message:        fmt.Sprintf("failed to insert line %d", row.LineNumber),
				Timestamp:      time.Now().UTC(),
			})
			return
		}

		processed++
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("failed to commit transaction", slog.String("error", err.Error()))
		tx.Rollback()
		logger.Error("transaction rolled back", slog.Bool("transaction_rolled_back", true))
		s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
			RequestID:      job.RequestID,
			Status:         "error",
			TotalRows:      totalRows,
			ProcessedRows:  0,
			DurationMS:     time.Since(startTime).Milliseconds(),
			ErrorCount:     1,
			FirstErrorLine: 0,
			Message:        "failed to commit transaction",
			Timestamp:      time.Now().UTC(),
		})
		return
	}

	logger.Info("csv processing completed",
		slog.Int("total_rows", totalRows),
		slog.Int("processed_rows", processed),
		slog.Int64("duration_ms", time.Since(startTime).Milliseconds()))

	s.publishFeedback(ctx, messaging.OpeningCSVFeedback{
		RequestID:      job.RequestID,
		Status:         "success",
		TotalRows:      totalRows,
		ProcessedRows:  processed,
		DurationMS:     time.Since(startTime).Milliseconds(),
		ErrorCount:     0,
		FirstErrorLine: 0,
		Message:        "csv processed successfully",
		Timestamp:      time.Now().UTC(),
	})
}

func (s *OpeningCSVService) publishFeedback(ctx context.Context, feedback messaging.OpeningCSVFeedback) {
	if s.producer == nil {
		s.logger.Error("feedback producer is not configured",
			slog.String("request_id", feedback.RequestID))
		return
	}

	if err := s.producer.PublishOpeningCSVFeedback(ctx, feedback); err != nil {
		s.logger.Error("failed to publish kafka feedback",
			slog.String("request_id", feedback.RequestID),
			slog.String("error", err.Error()))
	}
}
