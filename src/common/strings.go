package common

import (
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mohsen104/web-api/config"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	min := int(math.Pow(10, float64(config.GetConfig().Otp.Digits-1)))
	max := int(math.Pow(10, float64(config.GetConfig().Otp.Digits)))

	var num = rand.Intn(max-min) + min
	return strconv.Itoa(num)
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
