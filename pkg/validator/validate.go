package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func NewValidate(utrans *ut.UniversalTranslator, translatorMiddleware *middleware.TranslatorMiddleware) (*validator.Validate, error) {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	en := en.New()
	trans, found := utrans.GetTranslator(en.Locale())
	if !found {
		return nil, fmt.Errorf("translator not found: %v", en.Locale())
	}
	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	return v, nil
}
