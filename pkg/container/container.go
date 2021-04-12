package container

import (
	"reflect"
	"strings"

	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/asim/go-micro/v3/config"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewContainer(cfg config.Config) *dig.Container {
	container := dig.New()

	sqlDialect := cfg.Get("sql_dialect").String("")
	sqlDsn := cfg.Get("sql_dsn").String("")
	container.Provide(func() (*gorm.DB, error) {
		var dialector gorm.Dialector
		switch strings.ToLower(sqlDialect) {
		case "mysql":
			dialector = mysql.Open(sqlDsn)
		case "sqlite":
			dialector = sqlite.Open(sqlDsn)
		default:
			err := errors.New("invalid sql dialect")
			return nil, errors.WithStack(err)
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

	container.Provide(func(utrans *ut.UniversalTranslator) (*validator.Validate, error) {
		v := validator.New()
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			if jsonName != "-" {
				return jsonName
			}
			return field.Name
		})

		en := en.New()
		trans, found := utrans.GetTranslator(en.Locale())
		if !found {
			err := errors.Errorf("translator not found: %v", en.Locale())
			return nil, errors.WithStack(err)
		}

		err := en_translations.RegisterDefaultTranslations(v, trans)
		return v, errors.WithStack(err)
	})

	//--Middleware
	container.Provide(middleware.NewRecoverMiddleware)
	container.Provide(middleware.NewTransactionMiddleware)
	container.Provide(middleware.NewTranslatorMiddleware)
	container.Provide(middleware.NewErrorHandlerMiddleware)

	return container
}
