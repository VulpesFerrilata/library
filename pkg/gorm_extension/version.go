package model

import (
	extension_clause "github.com/VulpesFerrilata/library/pkg/gorm_extension/clause"
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
		extension_clause.NewVersionCreateClause(f),
	}
}

func (v version) UpdateClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		extension_clause.NewVersionUpdateClause(f),
	}
}

func (v version) DeleteClauses(f *schema.Field) []clause.Interface {
	return []clause.Interface{
		extension_clause.NewVersionDeleteClause(f),
	}
}
