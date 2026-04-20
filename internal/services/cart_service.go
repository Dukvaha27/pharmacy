package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrCartUserNotFound      = errors.New("user not found")
	ErrCartMedicineNotFound  = errors.New("medicine not found")
	ErrCartItemNotFound      = errors.New("cart item not found")
	ErrCartInvalidQuantity   = errors.New("quantity must be positive")
	ErrCartInsufficientStock = errors.New("insufficient stock")
)

type CartService interface {
	ClearCart(userID uint64) error
	GetByUserID(userID uint64) (*models.Cart, error)
	UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error
	DeleteItem(userID, itemID uint64) error
	AddItem(userID uint64, cartItemReq models.CartItemCreateRequest) error
}

type cartService struct {
	cartRepo     repository.CartRepository
	userRepo     repository.UserRepository
	medicineRepo repository.MedicineRepository
}

func NewCartService(
	cartRepo repository.CartRepository,
	userRepo repository.UserRepository,
	medicineRepo repository.MedicineRepository,
) CartService {
	return &cartService{
		cartRepo:     cartRepo,
		userRepo:     userRepo,
		medicineRepo: medicineRepo,
	}
}

func (s *cartService) ClearCart(userID uint64) error {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartUserNotFound
		}
		return err
	}

	return s.cartRepo.ClearCart(userID)
}

func (s *cartService) GetByUserID(userID uint64) (*models.Cart, error) {
	if _, err := s.userRepo.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartUserNotFound
		}
		return nil, err
	}

	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.Cart{
				UserID:     uint(userID),
				CartItems:  []models.CartItem{},
				TotalPrice: 0,
			}, nil
		}
		return nil, err
	}

	if cart.CartItems == nil {
		cart.CartItems = []models.CartItem{}
	}

	return cart, nil
}

func (s *cartService) AddItem(userID uint64, cartItemReq models.CartItemCreateRequest) error {
	if cartItemReq.Quantity == 0 {
		return ErrCartInvalidQuantity
	}

	if _, err := s.userRepo.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartUserNotFound
		}
		return err
	}

	medicine, err := s.medicineRepo.FindByID(cartItemReq.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartMedicineNotFound
		}
		return err
	}
	if medicine == nil {
		return ErrCartMedicineNotFound
	}

	if !medicine.InStock || cartItemReq.Quantity > uint64(medicine.StockQuantity) {
		return ErrCartInsufficientStock
	}

	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = &models.Cart{
				UserID:     uint(userID),
				CartItems:  []models.CartItem{},
				TotalPrice: 0,
			}
			if err := s.cartRepo.Create(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for _, existingItem := range cart.CartItems {
		if existingItem.MedicineID == cartItemReq.MedicineID {
			newQuantity := existingItem.Quantity + cartItemReq.Quantity
			if newQuantity > uint64(medicine.StockQuantity) {
				return ErrCartInsufficientStock
			}

			existingItem.Quantity = newQuantity
			existingItem.PricePerUnit = medicine.Price
			existingItem.LineTotal = newQuantity * medicine.Price

			if err := s.cartRepo.UpdateItem(&existingItem); err != nil {
				return err
			}

			newCartTotal := cart.TotalPrice + (cartItemReq.Quantity * medicine.Price)
			return s.cartRepo.SetCartTotalPrice(userID, newCartTotal)
		}
	}

	cartItem := &models.CartItem{
		MedicineID:   cartItemReq.MedicineID,
		Quantity:     cartItemReq.Quantity,
		LineTotal:    cartItemReq.Quantity * medicine.Price,
		PricePerUnit: medicine.Price,
		CartID:       cart.ID,
	}

	if err := s.cartRepo.AddItem(cartItem); err != nil {
		return err
	}

	newCartTotal := cart.TotalPrice + cartItem.LineTotal
	return s.cartRepo.SetCartTotalPrice(userID, newCartTotal)
}

func (s *cartService) UpdateItem(userID, itemID uint64, item *models.CartItemUpdateRequest) error {
	if item == nil || item.Quantity == nil || *item.Quantity == 0 {
		return ErrCartInvalidQuantity
	}

	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartUserNotFound
		}
		return err
	}

	var cartItem *models.CartItem
	for i := range cart.CartItems {
		if cart.CartItems[i].ID == uint(itemID) {
			cartItem = &cart.CartItems[i]
			break
		}
	}
	if cartItem == nil {
		return ErrCartItemNotFound
	}

	medicine, err := s.medicineRepo.FindByID(cartItem.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartMedicineNotFound
		}
		return err
	}
	if medicine == nil {
		return ErrCartMedicineNotFound
	}

	if !medicine.InStock || *item.Quantity > uint64(medicine.StockQuantity) {
		return ErrCartInsufficientStock
	}

	oldLineTotal := cartItem.LineTotal

	cartItem.Quantity = *item.Quantity
	cartItem.PricePerUnit = medicine.Price
	cartItem.LineTotal = (*item.Quantity) * medicine.Price

	if err := s.cartRepo.UpdateItem(cartItem); err != nil {
		return err
	}

	newCartTotal := cart.TotalPrice - oldLineTotal + cartItem.LineTotal
	return s.cartRepo.SetCartTotalPrice(userID, newCartTotal)
}

func (s *cartService) DeleteItem(userID, itemID uint64) error {
	cart, err := s.cartRepo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartUserNotFound
		}
		return err
	}

	var lineTotal uint64
	found := false
	for _, v := range cart.CartItems {
		if v.ID == uint(itemID) {
			lineTotal = v.LineTotal
			found = true
			break
		}
	}
	if !found {
		return ErrCartItemNotFound
	}

	if err := s.cartRepo.DeleteItem(userID, itemID); err != nil {
		return err
	}

	newCartTotal := cart.TotalPrice - lineTotal
	return s.cartRepo.SetCartTotalPrice(userID, newCartTotal)
}
