package init

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

func InitValidator(utrans *ut.UniversalTranslator) (*validator.Validate, error) {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonName == "-" || jsonName == "" {
			return field.Name
		}
		return jsonName
	})

	en := en.New()
	trans, found := utrans.GetTranslator(en.Locale())
	if !found {
		err := errors.Errorf("translator for given locale not found: %v", en.Locale())
		return nil, errors.WithStack(err)
	}

	err := en_translations.RegisterDefaultTranslations(v, trans)
	return v, errors.WithStack(err)
}
