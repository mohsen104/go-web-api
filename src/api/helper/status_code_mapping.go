package helper

import (
	"net/http"

	"github.com/mohsen104/web-api/pkg/service_errors"
)

var StatusCodeMapping = map[string]int{
	service_errors.OtpExists:   409,
	service_errors.OtpNotValid: 400,
	service_errors.OtpUsed:     409,
	service_errors.EmailExists: 409,
	service_errors.UsernameExists: 409,
}

func TranslateErrorToStatusCode(err error) int {
	if code, ok := StatusCodeMapping[err.Error()]; ok {
		return code
	}
	return http.StatusInternalServerError
}
