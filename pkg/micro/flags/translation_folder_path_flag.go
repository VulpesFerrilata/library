package flags

import "github.com/micro/cli/v2"

const translationFolderPath = "translation_folder_path"

var TranslationFolderPathFlag = &cli.StringFlag{
	Name:    translationFolderPath,
	EnvVars: []string{"MICRO_TRANSLATION_FOLDER_PATH"},
	Usage:   "Folder path contains translation information which will be imported by universal translator",
}

func GetTranslationFolderPath(ctx cli.Context) string {
	return ctx.String(translationFolderPath)
}
