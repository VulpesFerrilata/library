package clause

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type VersionCreateClause interface {
	clause.Interface
	gorm.StatementModifier
}

func NewVersionCreateClause(field *schema.Field) VersionCreateClause {
	return &versionCreateClause{
		field: field,
	}
}

type versionCreateClause struct {
	field *schema.Field
}

func (v versionCreateClause) Name() string {
	return ""
}

func (v versionCreateClause) Build(clause.Builder) {
}

func (v versionCreateClause) MergeClause(*clause.Clause) {
}

func (v versionCreateClause) ModifyStatement(stmt *gorm.Statement) {
	if c, ok := stmt.Clauses[clause.Values{}.Name()]; ok {
		if values, ok := c.Expression.(clause.Values); ok {
			for idx := 0; idx < len(values.Columns); idx++ {
				if values.Columns[idx].Name == v.field.DBName {
					for _, value := range values.Values {
						value[idx] = 1
					}
					break
				}
			}
		}
	}
}
