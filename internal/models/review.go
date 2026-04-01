package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID     uint `json:"user_id" gorm:"not null,index"`
	MedicineID uint `json:"medicine_id" gorm:"not null,index"`
	Medicine   Medicine
	Rating     uint   `json:"rating" binding:"required,gte=1,lte=5"`
	Text       string `json:"text" gorm:"not null"`
}

type ReviewCreateRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Rating     uint   `json:"rating" binding:"required,gte=1,lte=5"`
	Text       string `json:"text" binding:"required"`
}

type ReviewUpdateRequest struct {
	Rating     *uint   `json:"rating"`
	Text       *string `json:"text"`
}
