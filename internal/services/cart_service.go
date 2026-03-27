package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo repository.CartRepository
	userRepo repository.UserRepository
	medicineRepo repository.MedicineRepository
}

func NewCartService(cartRepo repository.CartRepository, userRepo repository.UserRepository, medmedicineRepo repository.MedicineRepository) CartService {
	return CartService{cartRepo: cartRepo, userRepo: userRepo, medicineRepo: medmedicineRepo}
}

func (s *CartService) ClearCart(userID uint64) error {
	return s.cartRepo.ClearCart(userID)
}

func (s *CartService) GetByUserID(userID uint64) (*models.Cart,error)  {
	cart, err := s.cartRepo.GetByUserID(userID)
	return cart, err
}

func (s *CartService) UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error {
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
				if err := s.cartRepo.UpdateCartTotalPrice(cart.UserID, sum); err != nil {
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

func (s *CartService) DeleteItem(userID, itemID uint) error {
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

func (s *CartService) AddItem(userID uint, cartItemReq models.CartItemCreateRequest) error {

	_, err := s.userRepo.GetByID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("User Not Found")
		} else {
			return err
		}
	}

	medicine,err :=s.medicineRepo.FindByID(uint(cartItemReq.MedicineID))
	if err!=nil {
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
