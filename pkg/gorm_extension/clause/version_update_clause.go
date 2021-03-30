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
			stmt.Clauses[clause.Set{}.Name()] = c
		}
	}

	if c, ok := stmt.Clauses[clause.Where{}.Name()]; ok {
		if where, ok := c.Expression.(clause.Where); ok && len(where.Exprs) > 1 {
			for _, expr := range where.Exprs {
				if orCond, ok := expr.(clause.OrConditions); ok && len(orCond.Exprs) == 1 {
					where.Exprs = []clause.Expression{
						clause.And(where.Exprs...),
					}

					c.Expression = where
					stmt.Clauses[clause.Where{}.Name()] = c
					break
				}
			}
		}
	}

	stmt.AddClause(clause.Where{Exprs: []clause.Expression{
		clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: v.field.DBName}, Value: nil},
	}})
}
