package repository

import (
	"errors"
	"pharmacy/internal/models"
	"strings"

	"gorm.io/gorm"
)

type MedicineRepository interface {
	Create(medicine *models.Medicine) error
	FindByID(id uint) (*models.Medicine, error)
	FindAll(filter models.MedicineFilter) ([]models.Medicine, int64, error)
	Update(medicine *models.Medicine) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
	Exists(id uint) (bool, error)
	UpdateAvgRating(id uint) error
}

type medicineRepository struct {
	db *gorm.DB
}

func NewMedicineRepository(db *gorm.DB) MedicineRepository {
	return &medicineRepository{db: db}
}

func (r *medicineRepository) Create(medicine *models.Medicine) error {
	if medicine.StockQuantity <= 0 {
		medicine.InStock = false
	} else {
		medicine.InStock = true
	}

	return r.db.Create(medicine).Error
}

func (r *medicineRepository) FindByID(id uint) (*models.Medicine, error) {
	var medicine models.Medicine
	err := r.db.
		Preload("Category").
		Preload("SubCategory").
		First(&medicine, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &medicine, nil
}

func (r *medicineRepository) FindAll(filter models.MedicineFilter) ([]models.Medicine, int64, error) {
	var medicines []models.Medicine
	var total int64

	query := r.db.Model(&models.Medicine{}).
		Preload("Category").
		Preload("SubCategory")

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.SubCategoryID != nil {
		query = query.Where("subcategory_id = ?", *filter.SubCategoryID)
	}

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}

	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}

	if filter.InStock != nil {
		query = query.Where("in_stock = ?", *filter.InStock)
	}

	if filter.PrescriptionRequired != nil {
		query = query.Where("prescription_required = ?", *filter.PrescriptionRequired)
	}

	if filter.MinRating != nil {
		query = query.Where("avg_rating >= ?", *filter.MinRating)
	}

	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(manufacturer) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortField := filter.SortBy
	if sortField == "" {
		sortField = "created_at"
	}

	allowedSortFields := map[string]bool{
		"name": true, "price": true, "avg_rating": true,
		"created_at": true, "updated_at": true, "stock_quantity": true,
	}

	if !allowedSortFields[sortField] {
		sortField = "created_at"
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	if filter.Page < 1 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.Limit

	err := query.
		Order(sortField + " " + sortOrder).
		Offset(offset).
		Limit(filter.Limit).
		Find(&medicines).Error

	return medicines, total, err
}

func (r *medicineRepository) Update(medicine *models.Medicine) error {
	if medicine.StockQuantity <= 0 {
		medicine.InStock = false
	} else {
		medicine.InStock = true
	}

	return r.db.Save(medicine).Error
}

func (r *medicineRepository) Delete(id uint) error {
	return r.db.Delete(&models.Medicine{}, id).Error
}

func (r *medicineRepository) UpdateStock(id uint, quantity int) error {
	result := r.db.Model(&models.Medicine{}).
		Where("id = ?", id).
		Update("stock_quantity", quantity)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return r.db.Model(&models.Medicine{}).
		Where("id = ?", id).
		Update("in_stock", quantity > 0).Error
}

func (r *medicineRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Medicine{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *medicineRepository) UpdateAvgRating(id uint) error {
	var avgRating float64

	err := r.db.
		Model(&models.Review{}).
		Where("medicine_id = ?", id).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating).Error
	if err != nil {
		return err
	}

	result := r.db.
		Model(&models.Medicine{}).
		Where("id = ?", id).
		Update("avg_rating", avgRating)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
