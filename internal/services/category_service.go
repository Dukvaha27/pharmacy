package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) Create(name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	existing, err := s.categoryRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	category := &models.Category{
		Name: name,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (s *CategoryService) GetByName(name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category, err := s.categoryRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	return category, nil
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	categories, err := s.categoryRepo.FindAll()
	if err != nil {
		return nil, err
	}

	if categories == nil {
		return []models.Category{}, nil
	}

	return categories, nil
}

func (s *CategoryService) Update(id uint, name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	existing, err := s.categoryRepo.FindByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.ID != id {
		return nil, errors.New("category with this name already exists")
	}

	category.Name = name
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Delete(id uint) error {
	exists, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if exists == nil {
		return errors.New("category not found")
	}

	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
