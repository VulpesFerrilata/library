package clause

import (
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
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
	if stmt.SQL.String() == "" {
		where, whereClauseExists := stmt.Clauses[clause.Where{}.Name()]

		set := callbacks.ConvertToAssignments(stmt)
		if len(set) == 0 {
			return
		}
		//remove version assignment value
		for idx, assignment := range set {
			if assignment.Column.Name == v.field.DBName {
				set = append(set[:idx], set[idx+1:]...)
			}
		}
		//assign version assignment value to auto increment
		assignment := clause.Assignment{
			Column: clause.Column{Name: v.field.DBName},
			Value:  gorm.Expr(v.field.DBName+"+ ?", 1),
		}
		set = append(set, assignment)

		//reset where clause
		if whereClauseExists {
			stmt.Clauses[clause.Where{}.Name()] = where
		} else {
			delete(stmt.Clauses, clause.Where{}.Name())
		}

		updatingValue := reflect.ValueOf(stmt.Dest)
		for updatingValue.Kind() == reflect.Ptr {
			updatingValue = updatingValue.Elem()
		}

		if !updatingValue.CanAddr() || stmt.Dest != stmt.Model {
			//add version field to criteria
			var criteriaFields = append(stmt.Schema.PrimaryFields, v.field)

			switch stmt.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				var criteriaExprs []clause.Expression
				for i := 0; i < stmt.ReflectValue.Len(); i++ {
					var exprs = make([]clause.Expression, len(criteriaFields))
					var notZero bool
					for idx, field := range criteriaFields {
						value, isZero := field.ValueOf(stmt.ReflectValue.Index(i))
						exprs[idx] = clause.Eq{Column: field.DBName, Value: value}
						notZero = notZero || !isZero
					}
					if notZero {
						criteriaExprs = append(criteriaExprs, clause.And(exprs...))
					}
				}
				stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.Or(criteriaExprs...)}})
			case reflect.Struct:
				for _, field := range criteriaFields {
					if value, isZero := field.ValueOf(stmt.ReflectValue); !isZero {
						stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: field.DBName, Value: value}}})
					}
				}
			}
		} else {
			if updatingValue.Kind() == reflect.Struct {
				for _, dbName := range stmt.Schema.DBNames {
					field := stmt.Schema.LookUpField(dbName)
					if field.PrimaryKey || field == v.field {
						if value, isZero := field.ValueOf(updatingValue); !isZero {
							stmt.AddClause(clause.Where{Exprs: []clause.Expression{clause.Eq{Column: field.DBName, Value: value}}})
						}
					}
				}
			}
		}

		stmt.AddClauseIfNotExists(clause.Update{})
		stmt.AddClause(set)
		stmt.Build("UPDATE", "SET", "WHERE")
	}
}
