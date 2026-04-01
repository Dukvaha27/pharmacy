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
}

func NewReviewService(reviewRepo repository.ReviewRepository, medicineRepo repository.MedicineRepository) reviewService {
	return reviewService{reviewRepo: reviewRepo, medicineRepo: medicineRepo}
}

func (s *reviewService) GetAll(medicineID uint64) (*[]models.Review, error) {
	if medicineID == 0 {
		return nil, errors.New("Invalid medicine ID")
	}

	reviews, err := s.reviewRepo.GetAll(medicineID)
	if err != nil {
		return nil, err
	}
	return &reviews, nil
}

func (s *reviewService) GetByID(reviewID uint64) (*models.Review, error) {
	if reviewID == 0 {
		return nil, errors.New("Invalid review ID")
	}

	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Review not found")
		} else {
			return nil, err
		}
	}
	return &review, nil
}

func (s *reviewService) Delete(reviewID uint64) error {
	if reviewID == 0 {
		return errors.New("Invalid review ID")
	}
	err := s.reviewRepo.Delete(reviewID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Review not found")
		} else {
			return err
		}
	}

	return nil
}
func (s *reviewService) Update(reviewID uint64, req models.ReviewUpdateRequest) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return errors.New("Review not found")
	}
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Text != nil {
		review.Text = *req.Text
	}

	return s.reviewRepo.Update(review)

}

func (s *reviewService) Create(req models.ReviewCreateRequest) error {
	medicine, err := s.medicineRepo.FindByID(req.MedicineID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Medicine Not Found")
		} else {
			return err
		}
	}

	review := models.Review{
		UserID:     req.UserID,
		MedicineID: req.MedicineID,
		Medicine:   *medicine,
		Rating:     req.Rating,
		Text:       req.Text,
	}

	return s.reviewRepo.Create(review)

}
