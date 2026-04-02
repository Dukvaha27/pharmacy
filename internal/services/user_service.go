package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
)

type UserService interface {
	GetByID(userID uint64) (*models.User, error)
	Update(userID uint64, user *models.UserUpdateRequest) error
	Create(user *models.UserCreateRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

// Create(user *models.User)error
// GetByID(userID uint64) (*models.User,error)
// Update(user *models.User) error

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(userID uint64) (*models.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) Update(userID uint64, user *models.UserUpdateRequest) error {
	oldUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user.DefaultAddress != nil {
		oldUser.DefaultAddress = user.DefaultAddress
	}
	if user.Email != nil {
		oldUser.Email = *user.Email
	}
	if user.FullName != nil {
		oldUser.FullName = *user.FullName
	}
	if user.Phone != nil {
		oldUser.Phone = *user.Phone
	}
	return s.userRepo.Update(oldUser)
}

func (s *userService) Create(user *models.UserCreateRequest) error {

	if len(user.Phone) != 11 {
		return errors.New("Номер телефона должен содержать 11 цифр!")
	}

	userModel := models.User{
		FullName:       user.FullName,
		Email:          user.Email,
		Phone:          user.Phone,
		DefaultAddress: &user.DefaultAddress,
	}

	return s.userRepo.Create(&userModel)
}
