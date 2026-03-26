package repository

import (
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(cart *models.Cart) error
	GetByUserID(userID uint64) (*models.Cart, error)

	AddItem(userID uint, cartItem *models.CartItem, cart *models.Cart) error
	UpdateCartTotalPrice(userID uint, summa int) error
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

func (r gormCartRepository) UpdateCartTotalPrice(userID uint, summa int) error {
	cart, err := r.GetByUserID(uint64(userID))
	if err != nil {
		return err
	}

	return r.db.Model(cart).Update("total_price", cart.TotalPrice+summa).Error
}

func (r gormCartRepository) ClearCart(userID uint64) error {
	cart, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.Model(cart).Association("CartItems").Clear()

}

func (r gormCartRepository) DeleteItem(cartID, itemID uint64) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}, itemID).Error
}

func (r gormCartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Model(&models.CartItem{}).Where("cart_id = ?", item.CartID).Select("*").Updates(item).Error
}

func (r gormCartRepository) Create(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r gormCartRepository) GetByUserID(userID uint64) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r gormCartRepository) AddItem(userID uint, cartItem *models.CartItem, cart *models.Cart) error {

	item := models.CartItem{
		MedicineID:   cartItem.MedicineID,
		Quantity:     cartItem.Quantity,
		CartID:       cart.ID,
		LineTotal:    cartItem.LineTotal,
		PricePerUnit: cartItem.PricePerUnit,
	}

	return r.db.Model(&cart).Association("CartItems").Append(&item)

}
