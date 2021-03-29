package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type VersionDeleteClause interface {
	clause.Interface
	gorm.StatementModifier
}

func NewVersionDeleteClause(field *schema.Field) VersionCreateClause {
	return &versionCreateClause{
		field: field,
	}
}

type versionDeleteClause struct {
	field *schema.Field
}

func (v versionDeleteClause) Name() string {
	return ""
}

func (v versionDeleteClause) Build(clause.Builder) {
}

func (v versionDeleteClause) MergeClause(*clause.Clause) {
}

func (v versionDeleteClause) ModifyStatement(stmt *gorm.Statement) {
}
