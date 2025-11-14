package validations

import (
	"errors"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/pkg/logging"
)

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	val, ok := fld.Field().Interface().(string)

	if !ok {
		return false
	}

	res, err := regexp.MatchString(`^09[0-9]{9}$`, val)

	if err != nil {
		logger := logging.NewLogger(config.GetConfig())
		logger.Error(logging.Validation, logging.MobileValidation, err.Error(), nil)
	}

	return res
}

type ValidationError struct {
	Property string `json:"property"`
	Tag      string `json:"tag"`
	Value    string `json:"value"`
	Message  string `json:"message"`
}

func GetValidationErrors(err error) *[]ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		for _, err := range err.(validator.ValidationErrors) {
			var el ValidationError
			el.Property = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			validationErrors = append(validationErrors, el)
		}
		return &validationErrors
	}

	return nil
}
