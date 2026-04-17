package repository

import (
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart) error
	GetByUserID(userID uint64) (*models.Cart, error)

	AddItem(cartItem *models.CartItem) error
	SetCartTotalPrice(userID uint64, total uint64) error
	UpdateItem(item *models.CartItem) error
	DeleteItem(userID, itemID uint64) error
	ClearCart(userID uint64) error
}

type gormCartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &gormCartRepository{db: db}
}

func (r *gormCartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *gormCartRepository) GetByUserID(userID uint64) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.
		Where("user_id = ?", userID).
		Preload("CartItems").
		First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *gormCartRepository) AddItem(cartItem *models.CartItem) error {
	return r.db.Create(cartItem).Error
}

func (r *gormCartRepository) SetCartTotalPrice(userID uint64, total uint64) error {
	cart, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.Model(&models.Cart{}).
		Where("id = ?", cart.ID).
		Update("total_price", total).Error
}

func (r *gormCartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Model(&models.CartItem{}).
		Where("id = ? AND cart_id = ?", item.ID, item.CartID).
		Updates(map[string]interface{}{
			"quantity":       item.Quantity,
			"line_total":     item.LineTotal,
			"price_per_unit": item.PricePerUnit,
		}).Error
}

func (r *gormCartRepository) DeleteItem(userID, itemID uint64) error {
	cart, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.
		Where("id = ? AND cart_id = ?", itemID, cart.ID).
		Delete(&models.CartItem{}).Error
}

func (r *gormCartRepository) ClearCart(userID uint64) error {
	cart, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("cart_id = ?", cart.ID).
			Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		return tx.Model(&models.Cart{}).
			Where("id = ?", cart.ID).
			Update("total_price", 0).Error
	})
}
