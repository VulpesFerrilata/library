package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type VersionUpdateClause interface {
	clause.Interface
	gorm.StatementModifier
}

func NewVersionUpdateClause(field *schema.Field) VersionUpdateClause {
	return &versionUpdateClause{
		field: field,
	}
}

type versionUpdateClause struct {
	field *schema.Field
}

func (v versionUpdateClause) Name() string {
	return ""
}

func (v versionUpdateClause) Build(clause.Builder) {
}

func (v versionUpdateClause) MergeClause(*clause.Clause) {
}

func (v versionUpdateClause) ModifyStatement(stmt *gorm.Statement) {
	if c, ok := stmt.Clauses[clause.Set{}.Name()]; ok {
		if set, ok := c.Expression.(clause.Set); ok {
			for idx, assignment := range set {
				if assignment.Column.Name == v.field.DBName {
					set = append(set[:idx], set[idx+1:]...)
				}
			}

			assignment := clause.Assignment{
				Column: clause.Column{Name: v.field.DBName},
				Value:  gorm.Expr(v.field.DBName+"+ ?", 1),
			}
			set = append(set, assignment)

			c.Expression = set
		}
	}
}
