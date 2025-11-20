package services

import (
	"context"

	"github.com/mohsen104/web-api/common"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/data/models"
	"github.com/mohsen104/web-api/pkg/logging"
	"github.com/mohsen104/web-api/pkg/service_errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	logger       logging.Logger
	cfg          *config.Config
	otpService   *OtpService
	tokenService *TokenService
	database     *gorm.DB
}

type GetOtpRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
}

type RegisterUserByUsername struct {
	FirstName string `json:"firstName" binding:"required,min=3"`
	LastName  string `json:"lastName" binding:"required,min=6"`
	Username  string `json:"username" binding:"required,min=5"`
	Email     string `json:"email" binding:"min=6,email"`
	Password  string `json:"password" binding:"required,password,min=6"`
}

type RegisterLoginByMobile struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Otp          string `json:"otp" binding:"required,min=6,max=6"`
}

type LoginByUsername struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=6"`
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

func (s *UserService) ExistsEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	if err := s.database.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").
		Where("email = ?", email).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) ExistsUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	if err := s.database.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").
		Where("username = ?", username).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) ExistsMobileNumber(ctx context.Context, mobileNumber string) (bool, error) {
	var exists bool
	if err := s.database.WithContext(ctx).Model(&models.User{}).
		Select("count(*) > 0").
		Where("mobile_number = ?", mobileNumber).
		Find(&exists).
		Error; err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}

func (s *UserService) GetDefaultRole(ctx context.Context) (roleId int, err error) {

	if err = s.database.WithContext(ctx).Model(&models.Role{}).
		Select("id").
		Where("name = ?", "default").
		First(&roleId).Error; err != nil {
		return 0, err
	}
	return roleId, nil
}

func (u *UserService) RegisterByUsername(ctx context.Context, req RegisterUserByUsername) error {
	user := models.User{Username: req.Username, FirstName: req.FirstName, LastName: req.LastName, Email: req.Email}

	exists, err := u.ExistsEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessages: service_errors.EmailExists}
	}
	exists, err = u.ExistsUsername(ctx, req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &service_errors.ServiceError{EndUserMessages: service_errors.UsernameExists}
	}

	bp := []byte(req.Password)
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return err
	}
	user.Password = string(hp)
	roleId, err := u.GetDefaultRole(ctx)
	if err != nil {
		u.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return err
	}

	tx := u.database.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: user.Id}).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return err
	}
	tx.Commit()
	return nil
}

func (u *UserService) RegisterAndLoginByMobileNumber(ctx context.Context, req RegisterLoginByMobile) (*TokenDetails, error) {
	err := u.otpService.ValidateOtp(req.MobileNumber, req.Otp)
	if err != nil {
		return nil, err
	}
	exists, err := u.ExistsMobileNumber(ctx, req.MobileNumber)
	if err != nil {
		return nil, err
	}

	user := models.User{MobileNumber: req.MobileNumber, Username: req.MobileNumber}

	if exists {
		var user models.User

		err = u.database.Model(&models.User{}).Where("username = ?", user.Username).Preload("UserRoles", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Role")
		}).Find(&user).Error
		if err != nil {
			return nil, err
		}

		tdto := TokenDto{UserId: user.Id, FirstName: user.FirstName, LastName: user.LastName, Email: user.Email, Username: user.Username}

		if len(*user.UserRoles) > 0 {
			for _, ur := range *user.UserRoles {
				tdto.Roles = append(tdto.Roles, ur.Role.Name)
			}
		}

		token, err := u.tokenService.GenerateToken(&tdto)

		if err != nil {
			return nil, err
		}

		return token, nil
	}

	bp := []byte(common.GenerateOtp())
	hp, err := bcrypt.GenerateFromPassword(bp, bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error(logging.General, logging.HashPassword, err.Error(), nil)
		return nil, err
	}
	user.Password = string(hp)
	roleId, err := u.GetDefaultRole(ctx)
	if err != nil {
		u.logger.Error(logging.Postgres, logging.DefaultRoleNotFound, err.Error(), nil)
		return nil, err
	}

	tx := u.database.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	err = tx.Create(&models.UserRole{RoleId: roleId, UserId: user.Id}).Error
	if err != nil {
		tx.Rollback()
		u.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return nil, err
	}
	tx.Commit()
	u.database.Model(&models.User{})
	return nil, nil
}
