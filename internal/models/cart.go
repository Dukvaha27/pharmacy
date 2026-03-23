package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     int        `json:"user_id" gorm:"not null"`
	CartItems  []CartItem `json:"cart_items"`
	TotalPrice int        `json:"total_price"`
}

type CartCreateRequest struct {
	UserID    int        `json:"user_id" binding:"required"`
	CartItems []CartItem `json:"cart_items" binding:"required"`
}

type CartUpdateRequest struct {
	UserID    *int        `json:"user_id"`
	CartItems *[]CartItem `json:"cart_items"`
}
