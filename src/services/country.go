package services

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/mohsen104/web-api/api/dto"
	"github.com/mohsen104/web-api/config"
	"github.com/mohsen104/web-api/data/db"
	"github.com/mohsen104/web-api/data/models"
	"github.com/mohsen104/web-api/pkg/logging"
	"gorm.io/gorm"
)

type CountryService struct {
	database *gorm.DB
	logger   logging.Logger
}

func NewCountryService(cfg *config.Config) *CountryService {
	return &CountryService{
		database: db.GetDb(),
		logger:   logging.NewLogger(cfg),
	}
}

// Create
func (s *CountryService) Create(ctx context.Context, req *dto.CreateUpdateCountryRequest) (*dto.CountryResponse, error) {
	country := models.Country{Name: req.Name}
	country.CreatedAt = time.Now().UTC()

	tx := s.database.WithContext(ctx).Begin()

	err := tx.Create(&country).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Insert, err.Error(), nil)
		return nil, err
	}

	tx.Commit()

	dto := &dto.CountryResponse{
		Id:   strconv.Itoa(country.Id),
		Name: country.Name,
	}

	return dto, nil
}

// Update
func (s *CountryService) Update(ctx context.Context, id int, req *dto.CreateUpdateCountryRequest) (*dto.CountryResponse, error) {
	updateMap := map[string]interface{}{
		"Name":        req.Name,
		"modified_at": &sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	tx := s.database.WithContext(ctx).Begin()

	err := tx.Model(&models.Country{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updateMap).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Update, err.Error(), nil)
		return nil, err
	}

	var country models.Country
	err = tx.Model(&models.Country{}).Where("id = ? AND deleted_at IS NULL", id).First(&country).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	tx.Commit()

	dto := &dto.CountryResponse{
		Id:   strconv.Itoa(country.Id),
		Name: country.Name,
	}

	return dto, nil
}

// Delete
func (s *CountryService) Delete(ctx context.Context, id int) error {
	tx := s.database.WithContext(ctx).Begin()

	deleteMap := map[string]interface{}{
		"deleted_at": &sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	err := tx.Model(&models.Country{}).Where("id = ? AND deleted_at IS NULL", id).Delete(deleteMap).Error

	if err != nil {
		tx.Rollback()
		s.logger.Error(logging.Postgres, logging.Delete, err.Error(), nil)
		return err
	}

	tx.Commit()
	return nil
}

// Get By Id
func (s *CountryService) GetById(ctx context.Context, id int) (*dto.CountryResponse, error) {
	country := &models.Country{}

	err := s.database.Where("id = ? AND deleted_at IS NULL", id).First(country).Error

	if err != nil {
		s.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return nil, err
	}

	dto := &dto.CountryResponse{
		Id:   strconv.Itoa(country.Id),
		Name: country.Name,
	}

	return dto, nil
}
