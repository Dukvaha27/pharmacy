package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	OrderID uint       `json:"order_id" gorm:"not null;index"`
	Amount  uint64     `json:"amount" gorm:"not null"`
	Status  string     `json:"status" gorm:"not null;default:'pending'"`
	Method  string     `json:"method" gorm:"not null"`
	PaidAt  *time.Time `json:"paid_at"`
}

type PaymentCreateRequest struct {
	Amount uint64 `json:"amount" binding:"required"`
	Method string `json:"method" binding:"required"`
}
