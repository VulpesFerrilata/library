package initialize

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

func InitValidator(utrans *ut.UniversalTranslator) (*validator.Validate, error) {
	v := validator.New()

	en := en.New()
	trans, found := utrans.GetTranslator(en.Locale())
	if !found {
		return nil, errors.Errorf("translator not found: %v", en.Locale())
	}

	err := en_translations.RegisterDefaultTranslations(v, trans)
	return v, errors.WithStack(err)
}
