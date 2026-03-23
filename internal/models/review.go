package models

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID     uint   `json:"user_id"`
	MedicineID uint   `json:"medicine_id"`
	Rating     uint   `json:"rating"`
	Text       string `json:"text"`
}
type ReviewCreateRequest struct {
	UserID     uint   `json:"user_id" binding:"required"`
	MedicineID uint   `json:"medicine_id" binding:"required"`
	Rating     uint   `json:"rating" binding:"required"`
	Text       string `json:"text" binding:"required"`
}
type ReviewUpdateRequest struct {
	UserID     *uint   `json:"user_id"`
	MedicineID *uint   `json:"medicine_id"`
	Rating     *uint   `json:"rating"`
	Text       *string `json:"text"`
}
