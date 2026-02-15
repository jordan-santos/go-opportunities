package repository

import (
	"opportunities/internal/schemas"

	"gorm.io/gorm"
)

type OpeningRepository interface {
	Create(opening *schemas.Openings) error
	Get(id string) (schemas.Openings, error)
	Delete(id string) error
	Update(opening *schemas.Openings) error
	List() ([]schemas.Openings, error)
}

type sqliteRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) OpeningRepository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Create(opening *schemas.Openings) error {
	return r.db.Create(opening).Error
}

func (r *sqliteRepository) Get(id string) (schemas.Openings, error) {
	var opening schemas.Openings
	if err := r.db.First(&opening, id).Error; err != nil {
		return schemas.Openings{}, err
	}
	return opening, nil
}

func (r *sqliteRepository) Delete(id string) error {
	return r.db.Delete(&schemas.Openings{}, id).Error
}

func (r *sqliteRepository) Update(opening *schemas.Openings) error {
	return r.db.Save(opening).Error
}

func (r *sqliteRepository) List() ([]schemas.Openings, error) {
	var openings []schemas.Openings
	if err := r.db.Find(&openings).Error; err != nil {
		return nil, err
	}
	return openings, nil
}
