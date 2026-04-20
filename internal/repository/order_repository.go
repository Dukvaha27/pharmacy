package repository

import (
	"errors"
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(*models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetByUserID(userID uint) ([]models.Order, error)
	UpdateStatus(id uint, status string) error
}

type gormOrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &gormOrderRepository{db: db}
}

func (r *gormOrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *gormOrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems").
		Preload("Payments").
		First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err

	}
	return &order, nil
}

func (r *gormOrderRepository) GetByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, err
}

func (r *gormOrderRepository) UpdateStatus(id uint, status string) error {
	result := r.db.
		Model(&models.Order{}).
		Where("id = ?", id).
		Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
