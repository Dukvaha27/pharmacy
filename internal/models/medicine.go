package models

import (
	"gorm.io/gorm"
)

type Medicine struct {
	gorm.Model
	Name                 string       `gorm:"type:varchar(255);not null;index" json:"name" validate:"required"`
	Description          string       `gorm:"type:text" json:"description"`
	Price                uint64       `gorm:"not null" json:"price" validate:"required,gt=0"`
	InStock              bool         `gorm:"default:true;index" json:"in_stock"`
	StockQuantity        int          `gorm:"default:0;not null" json:"stock_quantity" validate:"gte=0"`
	CategoryID           uint         `gorm:"not null;index" json:"category_id" validate:"required"`
	SubCategoryID        *uint        `gorm:"index" json:"subcategory_id"`
	Manufacturer         string       `gorm:"type:varchar(255)" json:"manufacturer"`
	PrescriptionRequired bool         `gorm:"default:false;index" json:"prescription_required"`
	AvgRating            float64      `gorm:"default:0;index" json:"avg_rating"`
	Category             Category     `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	SubCategory          *SubCategory `gorm:"foreignKey:SubCategoryID" json:"subcategory,omitempty"`
}

type MedicineCreateRequest struct {
	Name                 string `json:"name" binding:"required,min=1,max=255"`
	Description          string `json:"description"`
	Price                uint64 `json:"price" binding:"required,gt=0"`
	StockQuantity        int    `json:"stock_quantity" binding:"gte=0"`
	CategoryID           uint   `json:"category_id" binding:"required"`
	SubCategoryID        *uint  `json:"subcategory_id"`
	Manufacturer         string `json:"manufacturer"`
	PrescriptionRequired bool   `json:"prescription_required"`
}

type MedicineUpdateRequest struct {
	Name                 *string `json:"name" binding:"omitempty,min=1,max=255"`
	Description          *string `json:"description"`
	Price                *uint64 `json:"price" binding:"omitempty,gt=0"`
	StockQuantity        *int    `json:"stock_quantity" binding:"omitempty,gte=0"`
	CategoryID           *uint   `json:"category_id"`
	SubCategoryID        *uint   `json:"subcategory_id"`
	Manufacturer         *string `json:"manufacturer"`
	PrescriptionRequired *bool   `json:"prescription_required"`
}

type MedicineFilter struct {
	CategoryID           *uint    `form:"category_id"`
	SubCategoryID        *uint    `form:"subcategory_id"`
	MinPrice             *uint64  `form:"min_price"`
	MaxPrice             *uint64  `form:"max_price"`
	InStock              *bool    `form:"in_stock"`
	PrescriptionRequired *bool    `form:"prescription_required"`
	Search               *string  `form:"search"`
	MinRating            *float64 `form:"min_rating"`
	Page                 int      `form:"page,default=1"`
	Limit                int      `form:"limit,default=20"`
	SortBy               string   `form:"sort_by,default=created_at"`
	SortOrder            string   `form:"sort_order,default=desc"`
}
