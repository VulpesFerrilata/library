package init

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func NewValidate(utrans *ut.UniversalTranslator) (*validator.Validate, error) {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonName != "-" {
			return jsonName
		}
		return ""
	})

	en := en.New()
	trans, found := utrans.GetTranslator(en.Locale())
	if !found {
		return nil, errors.Errorf("translator not found: %v", en.Locale())
	}

	err := en_translations.RegisterDefaultTranslations(v, trans)
	return v, errors.WithStack(err)
}
