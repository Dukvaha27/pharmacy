package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	MedicineID   int  `json:"medicine_id"`
	Quantity     int  `json:"quantity"`
	LineTotal    int  `json:"line_total"`
	PricePerUnit int  `json:"price_per_unit"`
	CartID       uint `json:"cart_id"`
}

type CartItemCreateRequest struct {
	PricePerUnit int `json:"price_per_unit" binding:"required"`
	MedicineID   int `json:"medicine_id" binding:"required"`
	Quantity     int `json:"quantity" binding:"required"`
}

type CartItemUpdateRequest struct {
	Quantity *int `json:"quantity" binding:"omitempty,gte=0"`
}
