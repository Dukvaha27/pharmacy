package models

import "gorm.io/gorm"

type SubCategory struct {
	gorm.Model
	CategoryID uint      `gorm:"column:category_id;not null;index" json:"category_id" validate:"required"`
	Name       string   `gorm:"column:name;not null" json:"name" validate:"required"`
	Category   Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"-"`
}
