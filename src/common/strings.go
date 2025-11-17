package common

import (
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/mohsen104/web-api/config"
)

func GenerateOtp() string {
	rand.Seed(time.Now().UnixNano())
	min := int(math.Pow(10, float64(config.GetConfig().Otp.Digits-1)))
	max := int(math.Pow(10, float64(config.GetConfig().Otp.Digits)))

	var num = rand.Intn(max-min) + min
	return strconv.Itoa(num)
}
