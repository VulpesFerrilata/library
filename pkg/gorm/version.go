package gorm

import (
	clause_custom "github.com/VulpesFerrilata/library/pkg/gorm/clause"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Version interface {
	schema.CreateClausesInterface
	schema.UpdateClausesInterface
	schema.DeleteClausesInterface
}

func NewVersion(value int64) Version {
	return version(value)
}

type version int64

func (v version) CreateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionCreateClause(f),
	}
}

func (v version) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionUpdateClause(f),
	}
}

func (v version) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		clause_custom.NewVersionDeleteClause(f),
	}
}
