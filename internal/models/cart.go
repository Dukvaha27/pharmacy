package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     uint       `json:"user_id" gorm:"not null;uniqueIndex"`
	CartItems  []CartItem `json:"cart_items"`
	TotalPrice uint64     `json:"total_price"`
}

type CartItemsRequest struct {
	Items []CartItemCreateRequest `json:"items" binding:"required,min=1,dive"`
}
