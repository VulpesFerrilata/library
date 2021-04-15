package flags

import "github.com/micro/cli/v2"

const sqlDsn = "sql_dsn"

var SqlDsnFlag = &cli.StringFlag{
	Name:    sqlDsn,
	EnvVars: []string{"MICRO_SQL_DSN"},
	Usage:   "Connection string which used for sql dialect to initialize data source",
}

func GetSqlDsn(cli *cli.Context) string {
	return cli.String(sqlDsn)
}
