package utils

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
	trans     ut.Translator
}

var Validator *CustomValidator

func init() {
	// setup translations for nicer default error messages
	v := validator.New()
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, trans)
	Validator = &CustomValidator{
		validator: v,
		trans:     trans,
	}
}

type validationError struct {
	Field   string
	Message string
}

func (cv *CustomValidator) Validate(i interface{}) error {
	errors := []validationError{}
	if err := cv.validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			e := validationError{
				Field:   err.Field(),
				Message: err.Translate(cv.trans),
			}
			errors = append(errors, e)
		}
		return echo.NewHTTPError(http.StatusBadRequest, errors)
	}
	return nil
}
