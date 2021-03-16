package model

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var StaleVersionErr = errors.New("attempted to update a stale object")

type Model struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Version   int
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
