package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `gorm:"column:name;type:varchar(255);not null;uniqueIndex" json:"name" validate:"required,min=1,max=255"`
}
