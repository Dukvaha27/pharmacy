package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID          uint        `json:"user_id" gorm:"not null;index"`
	Status          string      `json:"status" gorm:"not null;default:'pending_payment'"`
	TotalPrice      uint64      `json:"total_price" gorm:"not null;default:0"`
	DiscountTotal   uint64      `json:"discount_total" gorm:"not null;default:0"`
	FinalPrice      uint64      `json:"final_price" gorm:"not null;default:0"`
	DeliveryAddress string      `json:"delivery_address" gorm:"not null"`
	Comment         string      `json:"comment"`
	OrderItems      []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
	Payments        []Payment   `json:"payments" gorm:"foreignKey:OrderID"`
}

type OrderCreateRequest struct {
	DeliveryAddress string  `json:"delivery_address" binding:"required"`
	Comment         *string `json:"comment"`
	Promocode       string  `json:"promocode"`
}

type OrderUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}
