package init

import (
	"github.com/VulpesFerrilata/library/config"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewTranslator(config *config.Config) (*ut.UniversalTranslator, error) {
	en := en.New()
	utrans := ut.New(en, en)

	if err := utrans.Import(ut.FormatJSON, config.TranslationFolderPath); err != nil {
		return nil, errors.WithStack(err)
	}

	return utrans, errors.WithStack(utrans.VerifyTranslations())
}
