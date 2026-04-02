package repository

import (
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	GetAll(medicineID uint64) ([]models.Review, error)
	GetByID(reviewID uint64) (models.Review, error)
	Create(review models.Review) error
	Update(review models.Review) error
	Delete(reviewID uint64) error
}

type gormReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &gormReviewRepository{db: db}
}

func (r gormReviewRepository) Delete(reviewID uint64) error {
	return r.db.Delete(&models.Review{}, reviewID).Error
}

func (r gormReviewRepository) Update(review models.Review) error {
	return r.db.Model(&review).Where("id = ?", review.ID).Updates(review).Error
}

func (r gormReviewRepository) Create(review models.Review) error {
	return r.db.Create(&review).Error
}

func (r gormReviewRepository) GetAll(medicineID uint64) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("medicine_id = ?", medicineID).Find(&reviews).Error
	return reviews, err
}

func (r gormReviewRepository) GetByID(reviewID uint64) (models.Review, error) {
	var review models.Review
	err := r.db.First(&review, reviewID).Error
	return review, err
}
