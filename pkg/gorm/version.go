package gorm

import (
	clause_custom "github.com/VulpesFerrilata/library/pkg/gorm/clause"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

var _ schema.CreateClausesInterface = new(Version)
var _ schema.UpdateClausesInterface = new(Version)
var _ schema.DeleteClausesInterface = new(Version)
var _ callbacks.AfterUpdateInterface = new(Version)
var _ callbacks.AfterDeleteInterface = new(Version)

type Version int64

func (v Version) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionCreateClause(f),
	}
}

func (v Version) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionUpdateClause(f),
	}
}

func (v Version) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionDeleteClause(f),
	}
}

func (v Version) AfterUpdate(tx *gorm.DB) error {
	if tx.RowsAffected == 0 {
		return StaleObjectErr
	}
	return nil
}

func (v Version) AfterDelete(tx *gorm.DB) error {
	if tx.RowsAffected == 0 {
		return StaleObjectErr
	}
	return nil
}
