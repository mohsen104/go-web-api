package services

import (
	"github.com/mohsen104/web-api/common"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/pkg/logging"
	"gorm.io/gorm"
)

type UserService struct {
	logger     logging.Logger
	cfg        *config.Config
	otpService *OtpService
	database   *gorm.DB
}

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
}

func NewUserService(cfg *config.Config) *UserService {
	database := db.GetDb()
	logger := logging.NewLogger(cfg)
	return &UserService{
		cfg:        cfg,
		database:   database,
		logger:     logger,
		otpService: NewOtpService(cfg),
	}
}

func (s *UserService) SendOtp(req *GetOtpRequest) error {
	otp := common.GenerateOtp()
	err := s.otpService.SetOtp(req.MobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}
