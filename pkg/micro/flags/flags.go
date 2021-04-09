package flags

import "github.com/micro/cli/v2"

var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    SqlDialect,
		EnvVars: []string{"MICRO_SQL_DIALECT"},
		Usage:   "Sql dialect for storing data, currently support sqlite and mysql",
	},
	&cli.StringFlag{
		Name:    SqlDsn,
		EnvVars: []string{"MICRO_SQL_DSN"},
		Usage:   "Connection string which used for sql dialect to initialize data source",
	},
	&cli.StringFlag{
		Name:    TranslationFolderPath,
		EnvVars: []string{"MICRO_TRANSLATION_FOLDER_PATH"},
		Usage:   "Folder path contains translation information which will be imported by universal translator",
	},
}
