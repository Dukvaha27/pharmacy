package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName       string `json:"full_name" gorm:"not null"`
	Email          string `json:"email" gorm:"not null"`
	Phone          string `json:"phone" gorm:"not null"`
	DefaultAddress string `json:"default_address"`
}

type UserCreateRequest struct {
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required,len=11"`
	DefaultAddress string `json:"default_address" binding:"min=5"`
}

type UserUpdateRequest struct {
	FullName       *string `json:"full_name"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	DefaultAddress *string `json:"default_address"`
}
