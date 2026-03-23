package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FullName      string `json:"full_name" gorm:"not null"`
	Email         string `json:"email" gorm:"not null"`
	Phone         string `json:"phone" gorm:"not null"`
	DefaultAdress string `json:"defaul_adress"`
}

type UserCreateRequest struct {
	FullName      string `json:"full_name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	DefaultAdress string `json:"defaul_adress"`
}

type UserUpdateRequest struct {
	FullName      *string `json:"full_name"`
	Email         *string `json:"email"`
	Phone         *string `json:"phone"`
	DefaultAdress *string `json:"defaul_adress"`
}
