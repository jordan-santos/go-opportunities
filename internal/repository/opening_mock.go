package repository

import (
	"opportunities/internal/schemas"

	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type OpeningRepositoryMock struct {
	mock.Mock
}

func (m *OpeningRepositoryMock) Create(opening *schemas.Openings) error {
	args := m.Called(opening)
	return args.Error(0)
}

func (m *OpeningRepositoryMock) CreateWithTx(tx *gorm.DB, opening *schemas.Openings) error {
	args := m.Called(tx, opening)
	return args.Error(0)
}

func (m *OpeningRepositoryMock) BeginTx() (*gorm.DB, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *OpeningRepositoryMock) Get(id string) (schemas.Openings, error) {
	args := m.Called(id)
	return args.Get(0).(schemas.Openings), args.Error(1)
}

func (m *OpeningRepositoryMock) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *OpeningRepositoryMock) Update(opening *schemas.Openings) error {
	args := m.Called(opening)
	return args.Error(0)
}

func (m *OpeningRepositoryMock) List() ([]schemas.Openings, error) {
	args := m.Called()
	return args.Get(0).([]schemas.Openings), args.Error(1)
}
