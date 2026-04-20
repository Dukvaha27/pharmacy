package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	MedicineID   uint   `json:"medicine_id"`
	Quantity     uint64 `json:"quantity"`
	LineTotal    uint64 `json:"line_total"`
	PricePerUnit uint64 `json:"price_per_unit"`
	CartID       uint   `json:"cart_id"`
}

type CartItemCreateRequest struct {
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Quantity   uint64 `json:"quantity" binding:"required,gt=0"`
}

type CartItemUpdateRequest struct {
	Quantity *uint64 `json:"quantity" binding:"omitempty,gt=0"`
}
