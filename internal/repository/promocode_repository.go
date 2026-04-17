package repository

import (
	"errors"
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type PromocodeRepository interface {
	Create(promocode *models.Promocode) error
	GetAll() ([]models.Promocode, error)
	GetByID(id uint) (*models.Promocode, error)
	GetByCode(code string) (*models.Promocode, error)
	Update(id uint, updates *models.PromocodeUpdateRequest) error
	Delete(id uint) error
}

type gormPromocodeRepository struct {
	db *gorm.DB
}

func NewPromocodeRepository(db *gorm.DB) PromocodeRepository {
	return &gormPromocodeRepository{db: db}
}

func (r *gormPromocodeRepository) Create(promocode *models.Promocode) error {
	return r.db.Create(promocode).Error
}

func (r *gormPromocodeRepository) GetAll() ([]models.Promocode, error) {
	var promocodes []models.Promocode
	err := r.db.Find(&promocodes).Error
	if err != nil {
		return nil, err
	}
	return promocodes, nil
}

func (r *gormPromocodeRepository) GetByID(id uint) (*models.Promocode, error) {
	var promocode models.Promocode
	err := r.db.First(&promocode, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &promocode, nil
}

func (r *gormPromocodeRepository) GetByCode(code string) (*models.Promocode, error) {
	var promocode models.Promocode
	err := r.db.Where("code = ?", code).First(&promocode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &promocode, nil
}

func (r *gormPromocodeRepository) Update(id uint, updates *models.PromocodeUpdateRequest) error {
	result := r.db.
		Model(&models.Promocode{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *gormPromocodeRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Promocode{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
