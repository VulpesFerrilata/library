package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/VulpesFerrilata/library/pkg/middleware"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

func NewValidate(utrans *ut.UniversalTranslator, translatorMiddleware *middleware.TranslatorMiddleware) (*validator.Validate, error) {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if jsonName != "" {
			return jsonName
		}

		gormTags := strings.Split(field.Tag.Get("gorm"), ";")
		for _, gormTag := range gormTags {
			fieldTags := strings.Split(gormTag, ":")
			if fieldTags[0] == "column" {
				return fieldTags[1]
			}
		}

		return ""
	})

	en := en.New()
	trans, found := utrans.GetTranslator(en.Locale())
	if !found {
		return nil, fmt.Errorf("translator not found: %v", en.Locale())
	}
	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, errors.Wrap(err, "validator.NewValidate")
	}

	return v, nil
}
