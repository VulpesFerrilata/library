package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var StaleVersionErr = errors.New("attempted to update a stale object")

type Model struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int `gorm:"index"`
}

func (m Model) beforeModify(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	tx.Statement.Where("version = ?", m.Version)
	return nil
}

func (m Model) afterModify(tx *gorm.DB) error {
	if tx.RowsAffected == 0 {
		currentVersion := 0
		err := tx.Model(tx.Statement.Model).Where("id = ?", m.ID).
		return StaleVersionErr
	}
	return nil
}

func (m Model) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", 1)
	return nil
}

func (m Model) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	tx.Statement.Where("version = ?", m.Version)
	return nil
}

func (m Model) AfterUpdate(tx *gorm.DB) error {
	if tx.RowsAffected == 0 {

		return StaleVersionErr
	}
	return nil
}

func (m Model) BeforeDelete(tx *gorm.DB) error {
	tx.Statement.SetColumn("version", gorm.Expr("version + 1"))
	tx.Statement.Where("version = ?", m.Version)
	return nil
}

func (m Model) AfterDelete(tx *gorm.DB) error {
	if tx.RowsAffected == 0 {

		return StaleVersionErr
	}
	return nil
}
