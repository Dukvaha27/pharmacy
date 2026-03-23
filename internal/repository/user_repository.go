package repository

import (
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User)error
	GetByID(id uint64) (*models.User,error)
	Update(user *models.User) error
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &gormUserRepository{db: db}
}

type gormUserRepository struct {
	db *gorm.DB
}

func (r *gormUserRepository) Create(user *models.User) error{
	return r.db.Create(user).Error
}

func (r *gormUserRepository) Update(user *models.User) error{
	return r.db.Model(&models.User{}).Where("id = ?", user.ID).Select("*").Updates(user).Error
}

func (r *gormUserRepository) GetByID(id uint64) (*models.User,error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err !=nil {
		return nil, err
	}
	return &user,nil
}

