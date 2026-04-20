package models

import (
	"time"

	"gorm.io/gorm"
)

type Promocode struct {
	gorm.Model
	Code           string    `json:"code" gorm:"not null;uniqueIndex"`
	Description    string    `json:"description"`
	DiscountType   string    `json:"discount_type" gorm:"not null"`
	DiscountValue  uint64    `json:"discount_value" gorm:"not null"`
	ValidFrom      time.Time `json:"valid_from"`
	ValidTo        time.Time `json:"valid_to"`
	MaxUses        int       `json:"max_uses" gorm:"default:0"`
	MaxUsesPerUser int       `json:"max_uses_per_user" gorm:"default:0"`
	UsedCount      int       `json:"used_count" gorm:"default:0"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`
}

type PromocodeCreateRequest struct {
	Code           string     `json:"code" binding:"required"`
	Description    *string    `json:"description"`
	DiscountType   string     `json:"discount_type" binding:"required,oneof=percent fixed"`
	DiscountValue  uint64     `json:"discount_value" binding:"required"`
	ValidFrom      *time.Time `json:"valid_from"`
	ValidTo        *time.Time `json:"valid_to"`
	MaxUses        *int       `json:"max_uses"`
	MaxUsesPerUser *int       `json:"max_uses_per_user"`
	IsActive       *bool      `json:"is_active"`
}

type PromocodeUpdateRequest struct {
	Code          *string `json:"code"`
	Description   *string `json:"description"`
	DiscountType  *string `json:"discount_type" binding:"oneof=percent fixed"`
	DiscountValue *uint64 `json:"discount_value"`
	UsedCount     *int    `json:"used_count"`
	IsActive      *bool   `json:"is_active"`
}

type PromocodeCheckRequest struct {
	Code        string `json:"code" binding:"required"`
	OrderAmount uint64 `json:"order_amount" binding:"required"`
}
