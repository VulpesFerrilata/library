package container

import (
	"strings"

	"github.com/asim/go-micro/v3/config"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewContainer(cfg config.Config) *dig.Container {
	container := dig.New()

	sqlDialect := cfg.Get("sql_dialect").String("sqlite")
	sqlDsn := cfg.Get("sql_dsn").String("./db.sqlite")
	container.Provide(func() (*gorm.DB, error) {
		var dialector gorm.Dialector
		switch strings.ToLower(sqlDialect) {
		case "mysql":
			dialector = mysql.Open(sqlDsn)
		case "sqlite":
			dialector = sqlite.Open(sqlDsn)
		default:
			return nil, errors.New("invalid sql name")
		}

		db, err := gorm.Open(dialector, &gorm.Config{})
		return db, errors.WithStack(err)
	})

	translationFolderPath := cfg.Get("translation_folder_path").String("./translation")
	container.Provide(func() (*ut.UniversalTranslator, error) {
		en := en.New()
		utrans := ut.New(en, en)

		if err := utrans.Import(ut.FormatJSON, translationFolderPath); err != nil {
			return nil, errors.WithStack(err)
		}

		return utrans, errors.WithStack(utrans.VerifyTranslations())
	})

	return container
}
