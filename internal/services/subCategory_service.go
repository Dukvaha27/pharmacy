package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
)

type SubCategoryService struct {
	subCategoryRepo repository.SubCategoryRepository
	categoryRepo    repository.CategoryRepository
}

func NewSubCategoryService(
	subCategoryRepo repository.SubCategoryRepository,
	categoryRepo repository.CategoryRepository,
) SubCategoryService {
	return SubCategoryService{
		subCategoryRepo: subCategoryRepo,
		categoryRepo:    categoryRepo,
	}
}

func (s *SubCategoryService) Create(categoryID uint, name string) (*models.SubCategory, error) {
	if name == "" {
		return nil, errors.New("subcategory name cannot be empty")
	}

	category, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	subCategories, err := s.subCategoryRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	for _, sub := range subCategories {
		if sub.Name == name {
			return nil, errors.New("subcategory with this name already exists in the category")
		}
	}

	subCategory := &models.SubCategory{
		CategoryID: categoryID,
		Name:       name,
	}

	if err := s.subCategoryRepo.Create(subCategory); err != nil {
		return nil, err
	}

	return subCategory, nil
}

func (s *SubCategoryService) GetByID(id uint) (*models.SubCategory, error) {
	subCategory, err := s.subCategoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if subCategory == nil {
		return nil, errors.New("subcategory not found")
	}

	return subCategory, nil
}

func (s *SubCategoryService) GetByName(name string) (*models.SubCategory, error) {
	if name == "" {
		return nil, errors.New("subcategory name cannot be empty")
	}

	subCategory, err := s.subCategoryRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if subCategory == nil {
		return nil, errors.New("subcategory not found")
	}

	return subCategory, nil
}

func (s *SubCategoryService) GetByCategoryID(categoryID uint) ([]models.SubCategory, error) {
	category, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	subCategories, err := s.subCategoryRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, err
	}

	if subCategories == nil {
		return []models.SubCategory{}, nil
	}

	return subCategories, nil
}

func (s *SubCategoryService) GetAll() ([]models.SubCategory, error) {
	subCategories, err := s.subCategoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	if subCategories == nil {
		return []models.SubCategory{}, nil
	}

	return subCategories, nil
}

func (s *SubCategoryService) Update(id uint, name string) (*models.SubCategory, error) {
	if name == "" {
		return nil, errors.New("subcategory name cannot be empty")
	}

	subCategory, err := s.subCategoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if subCategory == nil {
		return nil, errors.New("subcategory not found")
	}

	subCategories, err := s.subCategoryRepo.FindByCategoryID(subCategory.CategoryID)
	if err != nil {
		return nil, err
	}

	for _, sub := range subCategories {
		if sub.Name == name && sub.ID != id {
			return nil, errors.New("subcategory with this name already exists in the category")
		}
	}

	subCategory.Name = name
	if err := s.subCategoryRepo.Update(subCategory); err != nil {
		return nil, err
	}

	return subCategory, nil
}

func (s *SubCategoryService) Delete(id uint) error {
	exists, err := s.subCategoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if exists == nil {
		return errors.New("subcategory not found")
	}

	if err := s.subCategoryRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s *SubCategoryService) GetCategoryWithSubCategories(categoryID uint) (*models.Category, []models.SubCategory, error) {
	category, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, nil, err
	}
	if category == nil {
		return nil, nil, errors.New("category not found")
	}

	subCategories, err := s.subCategoryRepo.FindByCategoryID(categoryID)
	if err != nil {
		return nil, nil, err
	}

	return category, subCategories, nil
}
