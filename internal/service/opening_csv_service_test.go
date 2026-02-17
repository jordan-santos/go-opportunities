package service

import (
	"context"
	"errors"
	"testing"

	"opportunities/internal/messaging"
	"opportunities/internal/repository"
	"opportunities/internal/schemas"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type feedbackProducerSpy struct {
	messages []messaging.OpeningCSVFeedback
}

func (f *feedbackProducerSpy) PublishOpeningCSVFeedback(_ context.Context, feedback messaging.OpeningCSVFeedback) error {
	f.messages = append(f.messages, feedback)
	return nil
}

type failOnSecondInsertRepo struct {
	db      *gorm.DB
	inserts int
}

func (r *failOnSecondInsertRepo) Create(opening *schemas.Openings) error {
	return r.db.Create(opening).Error
}

func (r *failOnSecondInsertRepo) CreateWithTx(tx *gorm.DB, opening *schemas.Openings) error {
	r.inserts++
	if r.inserts == 2 {
		return errors.New("forced insert error")
	}
	return tx.Create(opening).Error
}

func (r *failOnSecondInsertRepo) BeginTx() (*gorm.DB, error) {
	tx := r.db.Begin()
	return tx, tx.Error
}

func (r *failOnSecondInsertRepo) Get(id string) (schemas.Openings, error) {
	var opening schemas.Openings
	err := r.db.First(&opening, id).Error
	return opening, err
}

func (r *failOnSecondInsertRepo) Delete(id string) error {
	return r.db.Delete(&schemas.Openings{}, id).Error
}

func (r *failOnSecondInsertRepo) Update(opening *schemas.Openings) error {
	return r.db.Save(opening).Error
}

func (r *failOnSecondInsertRepo) List() ([]schemas.Openings, error) {
	var openings []schemas.Openings
	err := r.db.Find(&openings).Error
	return openings, err
}

func TestOpeningCSVService_ProcessJobSuccess(t *testing.T) {
	db := openTestDB(t)
	repo := repository.New(db)
	producer := &feedbackProducerSpy{}
	svc := NewOpeningCSVService(repo, producer, 1)

	content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,2000\n")
	svc.processJob(context.Background(), OpeningCSVJob{
		RequestID: "req-success",
		Content:   content,
	})

	var openings []schemas.Openings
	if err := db.Find(&openings).Error; err != nil {
		t.Fatalf("unexpected db error: %v", err)
	}
	if len(openings) != 1 {
		t.Fatalf("expected 1 opening persisted, got %d", len(openings))
	}
	if len(producer.messages) != 1 {
		t.Fatalf("expected 1 feedback message, got %d", len(producer.messages))
	}
	if producer.messages[0].Status != "success" {
		t.Fatalf("expected success status, got %s", producer.messages[0].Status)
	}
	if producer.messages[0].ProcessedRows != 1 {
		t.Fatalf("expected processed_rows 1, got %d", producer.messages[0].ProcessedRows)
	}
}

func TestOpeningCSVService_ProcessJobValidationError(t *testing.T) {
	db := openTestDB(t)
	repo := repository.New(db)
	producer := &feedbackProducerSpy{}
	svc := NewOpeningCSVService(repo, producer, 1)

	content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,0\n")
	svc.processJob(context.Background(), OpeningCSVJob{
		RequestID: "req-validation",
		Content:   content,
	})

	var openings []schemas.Openings
	if err := db.Find(&openings).Error; err != nil {
		t.Fatalf("unexpected db error: %v", err)
	}
	if len(openings) != 0 {
		t.Fatalf("expected no openings persisted, got %d", len(openings))
	}
	if len(producer.messages) != 1 {
		t.Fatalf("expected 1 feedback message, got %d", len(producer.messages))
	}
	if producer.messages[0].Status != "error" {
		t.Fatalf("expected error status, got %s", producer.messages[0].Status)
	}
	if producer.messages[0].FirstErrorLine != 2 {
		t.Fatalf("expected first_error_line 2, got %d", producer.messages[0].FirstErrorLine)
	}
}

func TestOpeningCSVService_ProcessJobRollbackOnInsertError(t *testing.T) {
	db := openTestDB(t)
	repo := &failOnSecondInsertRepo{db: db}
	producer := &feedbackProducerSpy{}
	svc := NewOpeningCSVService(repo, producer, 1)

	content := []byte("role,company,location,remote,link,salary\nGo Dev,Acme,BR,true,https://acme.com,1000\nGo Dev 2,Acme,BR,false,https://acme2.com,1000\n")
	svc.processJob(context.Background(), OpeningCSVJob{
		RequestID: "req-rollback",
		Content:   content,
	})

	var openings []schemas.Openings
	if err := db.Find(&openings).Error; err != nil {
		t.Fatalf("unexpected db error: %v", err)
	}
	if len(openings) != 0 {
		t.Fatalf("expected rollback to persist 0 openings, got %d", len(openings))
	}
	if len(producer.messages) != 1 {
		t.Fatalf("expected 1 feedback message, got %d", len(producer.messages))
	}
	if producer.messages[0].Status != "error" {
		t.Fatalf("expected error status, got %s", producer.messages[0].Status)
	}
	if producer.messages[0].FirstErrorLine != 3 {
		t.Fatalf("expected first_error_line 3, got %d", producer.messages[0].FirstErrorLine)
	}
}

func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed opening test db: %v", err)
	}

	if err := db.AutoMigrate(&schemas.Openings{}); err != nil {
		t.Fatalf("failed migrating test db: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		t.Cleanup(func() {
			_ = sqlDB.Close()
		})
	}

	return db
}
