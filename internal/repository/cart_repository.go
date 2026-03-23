package repository

import "pharmacy/internal/models"

type CartRepository interface {
	GetByUserID(userID uint64) (*models.Cart,error)

	AddItem(c *models.CartItem) error
	
	UpdateItem(item *models.CartItem) error
	
	DeleteItem(itemID uint64) error
	ClearCart(userID uint64) error
}



