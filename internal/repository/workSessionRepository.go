package repository

import (
	"play-wails/internal/model"

	"github.com/google/uuid"
)

type WorkSessionRepository interface {
	Create(session *model.WorkSession) error
	FindByID(id uuid.UUID) (*model.WorkSession, error)
	Update(session *model.WorkSession) error
	ListByRunID(runID uuid.UUID) ([]*model.WorkSession, error)
	Delete(id uuid.UUID) error
}
