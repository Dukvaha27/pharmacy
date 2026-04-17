package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
	"time"

	"gorm.io/gorm"
)

var validTransitions = map[string][]string{
	"pending_payment": {"paid", "canceled"},
	"paid":            {"shipped"},
	"shipped":         {"completed"},
	"completed":       {},
	"canceled":        {},
}

type OrderService struct {
	OrderRepo    repository.OrderRepository
	CartRepo     repository.CartRepository
	PromoRepo    repository.PromocodeRepository
	MedicineRepo repository.MedicineRepository
}

func NewOrderService(
	orderRepo repository.OrderRepository,
	cartRepo repository.CartRepository,
	promoRepo repository.PromocodeRepository,
	medicineRepo repository.MedicineRepository,
) *OrderService {
	return &OrderService{
		OrderRepo:    orderRepo,
		CartRepo:     cartRepo,
		PromoRepo:    promoRepo,
		MedicineRepo: medicineRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint, req models.OrderCreateRequest) (*models.Order, error) {
	cart, err := s.CartRepo.GetByUserID(uint64(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("корзина не найдена")
		}
		return nil, err
	}

	if len(cart.CartItems) == 0 {
		return nil, errors.New("корзина пуста")
	}

	orderItems := make([]models.OrderItem, 0, len(cart.CartItems))
	for _, item := range cart.CartItems {
		medicine, err := s.MedicineRepo.FindByID(item.MedicineID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("лекарство не найдено")
			}
			return nil, err
		}
		if medicine == nil {
			return nil, errors.New("лекарство не найдено")
		}

		if !medicine.InStock || item.Quantity > uint64(medicine.StockQuantity) {
			return nil, errors.New("недостаточно товара на складе")
		}

		orderItems = append(orderItems, models.OrderItem{
			MedicineID:   item.MedicineID,
			MedicineName: medicine.Name,
			Quantity:     item.Quantity,
			PricePerUnit: item.PricePerUnit,
			LineTotal:    item.LineTotal,
		})
	}

	totalPrice := cart.TotalPrice
	var discountTotal uint64
	var usedPromo *models.Promocode

	if req.Promocode != "" {
		promo, err := s.PromoRepo.GetByCode(req.Promocode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("промокод не найден")
			}
			return nil, err
		}
		if promo == nil {
			return nil, errors.New("промокод не найден")
		}
		if !promo.IsActive {
			return nil, errors.New("промокод неактивен")
		}

		now := time.Now()
		if !promo.ValidFrom.IsZero() && now.Before(promo.ValidFrom) {
			return nil, errors.New("промокод ещё не действует")
		}
		if !promo.ValidTo.IsZero() && now.After(promo.ValidTo) {
			return nil, errors.New("срок действия промокода истёк")
		}
		if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
			return nil, errors.New("достигнут лимит использований промокода")
		}

		switch promo.DiscountType {
		case "percent":
			discountTotal = totalPrice * promo.DiscountValue / 100
		case "fixed":
			discountTotal = promo.DiscountValue
		default:
			return nil, errors.New("неизвестный тип скидки")
		}

		if discountTotal > totalPrice {
			discountTotal = totalPrice
		}

		usedPromo = promo
	}

	order := &models.Order{
		UserID:          userID,
		Status:          "pending_payment",
		TotalPrice:      totalPrice,
		DiscountTotal:   discountTotal,
		FinalPrice:      totalPrice - discountTotal,
		DeliveryAddress: req.DeliveryAddress,
		OrderItems:      orderItems,
	}

	if req.Comment != nil {
		order.Comment = *req.Comment
	}

	if err := s.OrderRepo.Create(order); err != nil {
		return nil, err
	}

	if err := s.CartRepo.ClearCart(uint64(userID)); err != nil {
		return nil, err
	}

	if usedPromo != nil {
		newUsedCount := usedPromo.UsedCount + 1
		if err := s.PromoRepo.Update(usedPromo.ID, &models.PromocodeUpdateRequest{
			UsedCount: &newUsedCount,
		}); err != nil {
			return nil, err
		}
	}

	return order, nil
}

func (s *OrderService) GetOrder(id uint) (*models.Order, error) {
	order, err := s.OrderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("заказ не найден")
		}
		return nil, err
	}
	if order == nil {
		return nil, errors.New("заказ не найден")
	}
	return order, nil
}

func (s *OrderService) GetByUserID(userID uint) ([]models.Order, error) {
	orders, err := s.OrderRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if orders == nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (s *OrderService) UpdateStatus(id uint, req models.OrderUpdateRequest) error {
	order, err := s.OrderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("заказ не найден")
		}
		return err
	}
	if order == nil {
		return errors.New("заказ не найден")
	}

	allowed, ok := validTransitions[order.Status]
	if !ok {
		return errors.New("неизвестный текущий статус заказа")
	}

	for _, next := range allowed {
		if next == req.Status {
			return s.OrderRepo.UpdateStatus(id, req.Status)
		}
	}

	return errors.New("недопустимый переход статуса: " + order.Status + " → " + req.Status)
}
