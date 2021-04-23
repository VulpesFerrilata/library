package gorm

import (
	"github.com/google/uuid"
)

type Model struct {
	ID uuid.UUID
	Version
}
