package gorm

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   Version
}
