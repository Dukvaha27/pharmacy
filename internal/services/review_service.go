package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

type ReviewService interface {
	Delete(reviewID uint64) error
	Update(reviewID uint64, req models.ReviewUpdateRequest) error
	Create(req models.ReviewCreateRequest) error
	GetAll(medicineID uint64) (*[]models.Review, error)
	GetByID(reviewID uint64) (*models.Review, error)
}

type reviewService struct {
	reviewRepo   repository.ReviewRepository
	medicineRepo repository.MedicineRepository
	orderRepo    repository.OrderRepository
}

func NewReviewService(
	reviewRepo repository.ReviewRepository,
	medicineRepo repository.MedicineRepository,
	orderRepo repository.OrderRepository,
) ReviewService {
	return &reviewService{
		reviewRepo:   reviewRepo,
		medicineRepo: medicineRepo,
		orderRepo:    orderRepo,
	}
}

func (s *reviewService) GetAll(medicineID uint64) (*[]models.Review, error) {
	if medicineID == 0 {
		return nil, errors.New("invalid medicine id")
	}

	reviews, err := s.reviewRepo.GetAll(medicineID)
	if err != nil {
		return nil, err
	}
	return &reviews, nil
}

func (s *reviewService) GetByID(reviewID uint64) (*models.Review, error) {
	if reviewID == 0 {
		return nil, errors.New("invalid review id")
	}

	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, err
	}

	return &review, nil
}

func (s *reviewService) Delete(reviewID uint64) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return err
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return err
	}

	return s.medicineRepo.UpdateAvgRating(review.MedicineID)
}

func (s *reviewService) Update(reviewID uint64, req models.ReviewUpdateRequest) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("review not found")
		}
		return err
	}

	if req.Rating != nil {
		if *req.Rating < 1 || *req.Rating > 5 {
			return errors.New("rating must be between 1 and 5")
		}
		review.Rating = *req.Rating
	}

	if req.Text != nil {
		review.Text = *req.Text
	}

	if err := s.reviewRepo.Update(review); err != nil {
		return err
	}

	return s.medicineRepo.UpdateAvgRating(review.MedicineID)
}

func (s *reviewService) Create(req models.ReviewCreateRequest) error {
	if req.Rating < 1 || req.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	medicine, err := s.medicineRepo.FindByID(req.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("medicine not found")
		}
		return err
	}
	if medicine == nil {
		return errors.New("medicine not found")
	}

	purchased, err := s.userPurchasedMedicine(req.UserID, req.MedicineID)
	if err != nil {
		return err
	}
	if !purchased {
		return errors.New("user has not purchased this medicine")
	}

	review := models.Review{
		UserID:     req.UserID,
		MedicineID: req.MedicineID,
		Rating:     req.Rating,
		Text:       req.Text,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return err
	}

	return s.medicineRepo.UpdateAvgRating(req.MedicineID)
}

func (s *reviewService) userPurchasedMedicine(userID, medicineID uint) (bool, error) {
	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		return false, err
	}

	for _, order := range orders {
		if order.Status == "canceled" || order.Status == "cancelled" {
			continue
		}

		fullOrder, err := s.orderRepo.GetByID(order.ID)
		if err != nil {
			return false, err
		}

		for _, item := range fullOrder.OrderItems {
			if item.MedicineID == medicineID {
				return true, nil
			}
		}
	}

	return false, nil
}
