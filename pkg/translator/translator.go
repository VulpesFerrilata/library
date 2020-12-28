package translator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
)

func NewTranslator() (*ut.UniversalTranslator, error) {
	en := en.New()
	utrans := ut.New(en, en)

	if err := utrans.Import(ut.FormatJSON, "translations"); err != nil {
		return nil, errors.Wrap(err, "translator.NewTranslator")
	}

	return utrans, utrans.VerifyTranslations()
}
