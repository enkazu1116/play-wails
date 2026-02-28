package repository

import (
	"play-wails/internal/model"

	"github.com/google/uuid"
)

type TimeRecordRepository interface {
	Create(record *model.TimeRecord) error
	FindByID(id uuid.UUID) (*model.TimeRecord, error)
	Update(record *model.TimeRecord) error
	List(excludeDeleted bool) ([]*model.TimeRecord, error)
	Delete(id uuid.UUID) error
}
