package services

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/constants"
	"github.com/mohsen104/web-api/data/cache"
	"github.com/mohsen104/web-api/pkg/logging"
	"github.com/mohsen104/web-api/pkg/service_errors"
)

type OtpService struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type OtpDto struct {
	Value string
	Used  bool
}

func NewOtpService(cfg *config.Config) *OtpService {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpService{logger: logger, cfg: cfg, redisClient: redis}
}

func (s *OtpService) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	val := &OtpDto{
		Value: otp,
		Used:  false,
	}
	res, err := cache.Get(s.redisClient, key, OtpDto{})
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessages: service_errors.OtpExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessages: service_errors.OtpUsed}
	}

	err = cache.Set(s.redisClient, key, val, s.cfg.Otp.ExpireTime*time.Second)

	if err != nil {
		return err
	}

	return nil
}

func (s *OtpService) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constants.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[OtpDto](s.redisClient, key, OtpDto{})

	if err != nil {
		return err
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessages: service_errors.OtpUsed}
	} else if err == nil && !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessages: service_errors.OtpNotValid}
	} else if err == nil && !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(s.redisClient, key, res, s.cfg.Otp.ExpireTime*time.Second)
		return err
	}
	return nil
}
