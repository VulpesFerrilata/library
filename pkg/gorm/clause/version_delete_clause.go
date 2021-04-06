package clause

import (
	"reflect"

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
	if stmt.SQL.String() == "" {
		if stmt.Schema != nil {
			criteriaFields := append(stmt.Schema.PrimaryFields, v.field)

			_, queryValues := schema.GetIdentityFieldValuesMap(stmt.ReflectValue, criteriaFields)
			column, values := schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if stmt.ReflectValue.CanAddr() && stmt.Dest != stmt.Model && stmt.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(stmt.Model), criteriaFields)
				column, values = schema.ToQueryValues(stmt.Table, stmt.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		stmt.AddClauseIfNotExists(clause.Delete{})
		stmt.AddClauseIfNotExists(clause.From{})
		stmt.Build("DELETE", "FROM", "WHERE")
	}
}
