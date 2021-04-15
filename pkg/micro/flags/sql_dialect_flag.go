package flags

import (
	"github.com/VulpesFerrilata/library/pkg/micro/flags/generic"
	"github.com/micro/cli/v2"
)

const sqlDialect = "sql_dialect"

func NewSqlDialectFlag() cli.Flag {
	return &cli.GenericFlag{
		Name:    sqlDialect,
		Value:   generic.NewStringGeneric("mysql", "sqlite"),
		EnvVars: []string{"MICRO_SQL_DIALECT"},
		Usage:   "Sql dialect for storing data, currently support sqlite and mysql",
	}
}

func GetSqlDialect(ctx *cli.Context) string {
	return ctx.String(sqlDialect)
}
