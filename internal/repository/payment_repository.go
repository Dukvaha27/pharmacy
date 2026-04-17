package repository

import (
	"errors"
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	GetByID(id uint) (*models.Payment, error)
	GetByOrderID(orderID uint) ([]models.Payment, error)
}

type gormPaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &gormPaymentRepository{db: db}
}

func (r *gormPaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *gormPaymentRepository) GetByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *gormPaymentRepository) GetByOrderID(orderID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.
		Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}
