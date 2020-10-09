package translator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

func NewTranslator() (*ut.UniversalTranslator, error) {
	en := en.New()
	utrans := ut.New(en, en)

	if err := utrans.Import(ut.FormatJSON, "translations"); err != nil {
		return utrans, err
	}

	return utrans, utrans.VerifyTranslations()
}
