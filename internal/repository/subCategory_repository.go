package repository

import (
	"errors"
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type SubCategoryRepository interface {
	Create(subCategory *models.SubCategory) error
	FindByID(id uint) (*models.SubCategory, error)
	FindByName(name string) (*models.SubCategory, error)
	FindByCategoryID(categoryID uint) ([]models.SubCategory, error)
	FindAll() ([]models.SubCategory, error)
	Update(subCategory *models.SubCategory) error
	Delete(id uint) error
}

type subCategoryRepository struct {
	db *gorm.DB
}

func NewSubCategoryRepository(db *gorm.DB) SubCategoryRepository {
	return &subCategoryRepository{db: db}
}

func (r *subCategoryRepository) Create(subCategory *models.SubCategory) error {
	return r.db.Create(subCategory).Error
}

func (r *subCategoryRepository) FindByID(id uint) (*models.SubCategory, error) {
	var subCategory models.SubCategory
	err := r.db.Preload("Category").First(&subCategory, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &subCategory, nil
}

func (r *subCategoryRepository) FindByName(name string) (*models.SubCategory, error) {
	var subCategory models.SubCategory
	err := r.db.Where("name = ?", name).First(&subCategory).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &subCategory, nil
}

// FindByCategoryID - без пагинации
func (r *subCategoryRepository) FindByCategoryID(categoryID uint) ([]models.SubCategory, error) {
	var subCategories []models.SubCategory
	err := r.db.Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Find(&subCategories).Error
	return subCategories, err
}

func (r *subCategoryRepository) FindAll() ([]models.SubCategory, error) {
	var subCategories []models.SubCategory

	err := r.db.Preload("Category").
		Order("created_at DESC").
		Find(&subCategories).Error

	return subCategories, err
}

func (r *subCategoryRepository) Update(subCategory *models.SubCategory) error {
	return r.db.Save(subCategory).Error
}

// Delete - удаление подкатегории
func (r *subCategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.SubCategory{}, id).Error
}
