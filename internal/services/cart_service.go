package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo       repository.CartRepository
	userRepository repository.UserRepository
}

func NewCartService(cartRepository repository.CartRepository) CartService {
	return CartService{cartRepo: cartRepository}
}


func (s *CartService) ClearCart(userID uint64) error {
	return s.cartRepo.ClearCart(userID)
}

func (s *CartService) UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	var cartItem models.CartItem
	var sum int
	var HasItemId bool
	for _, v := range cart.CartItems {
		if v.ID == uint(itemID) {
			HasItemId = true
			cartItem = v
			if item.Quantity != nil {
				sum = (*item.Quantity + cartItem.PricePerUnit) - cartItem.LineTotal
				cartItem.Quantity = *item.Quantity
			}

			sum = cartItem.PricePerUnit * cartItem.Quantity
			s.cartRepo.UpdateCartTotalPrice(cart.UserID, sum)
			break
		}
	}
	if HasItemId == false {
		return errors.New("Not Found Item ID")
	}
	return s.cartRepo.UpdateItem(&cartItem)
}

func (s *CartService) DeleteItem(userID, itemID uint) error {
	cart, err := s.cartRepo.GetByUserID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User Not Found")
		} else {
			return err
		}

	}

	var number int
	for _, v := range cart.CartItems {
		if v.ID == itemID {
			number = v.LineTotal
		}
	}

	s.cartRepo.UpdateCartTotalPrice(userID, -(number))

	return s.cartRepo.DeleteItem(uint64(cart.ID), uint64(itemID))
}

func (s *CartService) AddItem(userID uint, cartItemReq models.CartItemCreateRequest) error {

	

	_, errUser := s.userRepository.GetByID(uint64(userID))
	if errUser != nil {
		if errors.Is(errUser, gorm.ErrRecordNotFound) {
			return errors.New("User Not Found")
		}
	}

	cart, errCart := s.cartRepo.GetByUserID(uint64(userID))
	if errCart != nil {
		if errors.Is(errCart, gorm.ErrRecordNotFound) {
			cart = &models.Cart{
				UserID: userID,
			}

			if err := s.cartRepo.Create(cart); err != nil {
				return err
			}
		} else {
			return errCart
		}
	}
	cartItem := &models.CartItem{
		MedicineID:   cartItemReq.MedicineID,
		Quantity:     cartItemReq.Quantity,
		LineTotal:    cartItemReq.Quantity * cartItemReq.PricePerUnit,
		PricePerUnit: cartItemReq.PricePerUnit,
		CartID:       cart.ID,
	}

	s.cartRepo.UpdateCartTotalPrice(userID, cartItem.LineTotal)

	return s.cartRepo.AddItem(userID, cartItem, cart)
}
