package services

import (
	"pharmacy/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}


	// Create(user *models.User)error
	// GetByID(id uint64) (*models.User,error)
	// Update(user *models.User) error


// func (r *UserService) Create(req models.UserCreateRequest) (*models.User,error) {

// }