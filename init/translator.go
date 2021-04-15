package init

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func InitTranslator(translationFolderPath string) (*ut.UniversalTranslator, error) {
	en := en.New()
	utrans := ut.New(en, en)

	if err := utrans.Import(ut.FormatJSON, translationFolderPath); err != nil {
		return nil, errors.WithStack(err)
	}

	return utrans, errors.WithStack(utrans.VerifyTranslations())
}
