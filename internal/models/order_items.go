package models

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderID      uint   `json:"order_id" gorm:"not null;index"`
	MedicineID   uint   `json:"medicine_id" gorm:"not null;index"`
	MedicineName string `json:"medicine_name" gorm:"not null"`
	Quantity     uint64 `json:"quantity" gorm:"not null"`
	PricePerUnit uint64 `json:"price_per_unit" gorm:"not null"`
	LineTotal    uint64 `json:"line_total" gorm:"not null"`
}
