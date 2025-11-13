package validations

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func IranianMobileNumberValidator(fld validator.FieldLevel) bool {
	val, ok := fld.Field().Interface().(string)

	if !ok {
		return false
	}

	res, err := regexp.MatchString(`^09[0-9]{9}$`, val)

	if err!=nil {
		log.Print(err.Error())
	}

	return res
}
