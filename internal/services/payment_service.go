package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
	"time"
)

type PaymentService struct {
	PaymentRepo repository.PaymentRepository
	OrderRepo   repository.OrderRepository
}

func NewPaymentService(
	paymentRepo repository.PaymentRepository,
	orderRepo repository.OrderRepository,
) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
		OrderRepo:   orderRepo,
	}
}

func (s *PaymentService) Create(
	orderID,
	userID uint,
	req models.PaymentCreateRequest,
) (*models.Payment, *models.Order, error) {
	if req.Amount == 0 {
		return nil, nil, errors.New("amount must be greater than 0")
	}

	order, err := s.OrderRepo.GetByID(orderID)
	if err != nil {
		return nil, nil, err
	}
	if order == nil {
		return nil, nil, errors.New("заказ не найден")
	}
	if order.UserID != userID {
		return nil, nil, errors.New("нет доступа к этому заказу")
	}

	if order.Status == "cancelled" ||
		order.Status == "canceled" ||
		order.Status == "completed" ||
		order.Status == "delivered" {
		return nil, nil, errors.New("невозможно создать платёж для заказа со статусом: " + order.Status)
	}

	payments, err := s.PaymentRepo.GetByOrderID(orderID)
	if err != nil {
		return nil, nil, err
	}

	var totalSuccessPaid uint64
	for _, p := range payments {
		if p.Status == "success" {
			totalSuccessPaid += p.Amount
		}
	}

	if totalSuccessPaid >= order.FinalPrice {
		return nil, nil, errors.New("заказ уже полностью оплачен")
	}

	if totalSuccessPaid+req.Amount > order.FinalPrice {
		return nil, nil, errors.New("сумма успешных платежей не может превышать final_price заказа")
	}

	now := time.Now()
	payment := &models.Payment{
		OrderID: orderID,
		Amount:  req.Amount,
		Status:  "success",
		Method:  req.Method,
		PaidAt:  &now,
	}

	if err := s.PaymentRepo.Create(payment); err != nil {
		return nil, nil, err
	}

	totalSuccessPaid += req.Amount
	if totalSuccessPaid == order.FinalPrice && order.Status == "pending_payment" {
		if err := s.OrderRepo.UpdateStatus(orderID, "paid"); err != nil {
			return nil, nil, err
		}
		order.Status = "paid"
	}

	order.Payments = append(order.Payments, *payment)
	return payment, order, nil
}

func (s *PaymentService) GetByOrderID(orderID, userID uint) ([]models.Payment, error) {
	order, err := s.OrderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("заказ не найден")
	}
	if order.UserID != userID {
		return nil, errors.New("нет доступа к заказу")
	}

	return s.PaymentRepo.GetByOrderID(orderID)
}

func (s *PaymentService) GetByID(id uint) (*models.Payment, error) {
	payment, err := s.PaymentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.New("платеж не найден")
	}
	return payment, nil
}
