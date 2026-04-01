package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

type CartService interface {
	ClearCart(userID uint64) error 
	GetByUserID(userID uint64) (*models.Cart, error) 
	 UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error
	 DeleteItem(userID, itemID uint) error 
	 AddItem(userID uint, cartItemReq models.CartItemCreateRequest) error 
}

type cartService struct {
	cartRepo     repository.CartRepository
	userRepo     repository.UserRepository
	medicineRepo repository.MedicineRepository
}



func NewCartService(cartRepo repository.CartRepository, userRepo repository.UserRepository, medicineRepo repository.MedicineRepository) cartService {
	return cartService{cartRepo: cartRepo, userRepo: userRepo, medicineRepo: medicineRepo}
}

func (s *cartService) ClearCart(userID uint64) error {
	return s.cartRepo.ClearCart(userID)
}

func (s *cartService) GetByUserID(userID uint64) (*models.Cart, error) {
	cart, err := s.cartRepo.GetByUserID(userID)
	return cart, err
}

func (s *cartService) UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		return err
	}
	var cartItem models.CartItem
	var sum int
	var hasItemID bool
	for _, v := range cart.CartItems {
		if v.ID == uint(itemID) {
			hasItemID = true
			cartItem = v
			if item.Quantity != nil {
				sum = (*item.Quantity * cartItem.PricePerUnit) - cartItem.LineTotal
				if err := s.cartRepo.UpdateCartTotalPrice(cart.UserID, sum); err != nil { //    
					return err
				}
				cartItem.Quantity = *item.Quantity
				cartItem.LineTotal = *item.Quantity * cartItem.PricePerUnit
			}

			break
		}
	}
	if !hasItemID {
		return errors.New("Not Found Item ID")
	}
	return s.cartRepo.UpdateItem(&cartItem)
}

func (s *cartService) DeleteItem(userID, itemID uint) error {
	cart, err := s.cartRepo.GetByUserID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User Not Found")
		} else {
			return err
		}

	}

	var lineTotal int
	var hasItemID bool
	for _, v := range cart.CartItems {
		if v.ID == itemID {
			lineTotal = v.LineTotal
			hasItemID = true
			break
		}
	}
	if !hasItemID {
		return errors.New("ItemId not Found")
	}

	if err := s.cartRepo.UpdateCartTotalPrice(userID, -(lineTotal)); err != nil {
		return err
	}

	return s.cartRepo.DeleteItem(uint64(cart.ID), uint64(itemID))
}

func (s *cartService) AddItem(userID uint, cartItemReq models.CartItemCreateRequest) error {

	_, err := s.userRepo.GetByID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User Not Found")
		} else {
			return err
		}
	}

	medicine, err := s.medicineRepo.FindByID(uint(cartItemReq.MedicineID))
	if err != nil {
		return errors.New("Medicine not found")
	}

	cart, err := s.cartRepo.GetByUserID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &models.Cart{
				UserID: userID,
			}

			if err := s.cartRepo.Create(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	cartItem := &models.CartItem{
		MedicineID:   cartItemReq.MedicineID,
		Quantity:     cartItemReq.Quantity,
		LineTotal:    cartItemReq.Quantity * int(medicine.Price),
		PricePerUnit: int(medicine.Price),
		CartID:       cart.ID,
	}

	if err := s.cartRepo.UpdateCartTotalPrice(userID, cartItem.LineTotal); err != nil {
		return err
	}

	return s.cartRepo.AddItem(userID, cartItem, cart)
}
