package services

import (
	"errors"
	"net/mail"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

var phoneRegex = regexp.MustCompile(`^\+?\d{11,15}$`)

type UserService interface {
	GetByID(userID uint64) (*models.User, error)
	Update(userID uint64, user *models.UserUpdateRequest) error
	Create(user *models.UserCreateRequest) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetByID(userID uint64) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Update(userID uint64, user *models.UserUpdateRequest) error {
	if user == nil {
		return errors.New("empty update request")
	}

	oldUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if user.DefaultAddress != nil {
		trimmed := strings.TrimSpace(*user.DefaultAddress)
		oldUser.DefaultAddress = &trimmed
	}
	if user.Email != nil {
		email := strings.TrimSpace(*user.Email)
		if _, err := mail.ParseAddress(email); err != nil {
			return errors.New("invalid email")
		}
		oldUser.Email = email
	}
	if user.FullName != nil {
		oldUser.FullName = strings.TrimSpace(*user.FullName)
	}
	if user.Phone != nil {
		phone := strings.TrimSpace(*user.Phone)
		if !phoneRegex.MatchString(phone) {
			return errors.New("invalid phone format")
		}
		oldUser.Phone = phone
	}

	return s.userRepo.Update(oldUser)
}

func (s *userService) Create(user *models.UserCreateRequest) error {
	if user == nil {
		return errors.New("empty create request")
	}

	email := strings.TrimSpace(user.Email)
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("invalid email")
	}

	phone := strings.TrimSpace(user.Phone)
	if !phoneRegex.MatchString(phone) {
		return errors.New("invalid phone format")
	}

	fullName := strings.TrimSpace(user.FullName)
	address := strings.TrimSpace(user.DefaultAddress)

	userModel := models.User{
		FullName: fullName,
		Email:    email,
		Phone:    phone,
	}

	if address != "" {
		userModel.DefaultAddress = &address
	}

	return s.userRepo.Create(&userModel)
}
