package services

import (
	"errors"
	"fmt"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"
)

var (
	ErrMedicineNotFound    = errors.New("medicine not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrCategoryNotFound    = errors.New("category not found")
	ErrSubCategoryNotFound = errors.New("SubCategory not found")
	ErrInvalidSubCategory  = errors.New("subcategory does not belong to the specified category")
	ErrInvalidID           = errors.New("Invalid id")
)

type MedicineService interface {
	CreateMedicine(req models.MedicineCreateRequest) (*models.Medicine, error)
	GetMedicineByID(id uint) (*models.Medicine, error)
	ListMedicines(filter models.MedicineFilter) ([]models.Medicine, int64, error)
	UpdateMedicine(id uint, req models.MedicineUpdateRequest) (*models.Medicine, error)
	DeleteMedicine(id uint) error
	CheckStock(medicineID uint, quantity int) (*models.Medicine, error)
	ReserveStock(medicineID uint, quantity int) error
	ReleaseStock(medicineID uint, quantity int) error
}

type medicineService struct {
	medicineRepo    repository.MedicineRepository
	categoryRepo    repository.CategoryRepository
	subCategoryRepo repository.SubCategoryRepository
}

func NewMedicineService(
	medicineRepo repository.MedicineRepository,
	categoryRepo repository.CategoryRepository,
	subCategoryRepo repository.SubCategoryRepository,
) MedicineService {
	return &medicineService{
		medicineRepo:    medicineRepo,
		categoryRepo:    categoryRepo,
		subCategoryRepo: subCategoryRepo,
	}
}

func (s *medicineService) CreateMedicine(req models.MedicineCreateRequest) (*models.Medicine, error) {
	category, err := s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.SubCategoryID != nil {
		subCategory, err := s.subCategoryRepo.FindByID(*req.SubCategoryID)
		if err != nil {
			return nil, err
		}
		if subCategory == nil {
			return nil, ErrSubCategoryNotFound
		}

		if subCategory.CategoryID != req.CategoryID {
			return nil, ErrInvalidSubCategory
		}
	}

	medicine := &models.Medicine{
		Name:                 req.Name,
		Description:          req.Description,
		Price:                req.Price,
		StockQuantity:        req.StockQuantity,
		CategoryID:           req.CategoryID,
		SubCategoryID:        req.SubCategoryID,
		Manufacturer:         req.Manufacturer,
		PrescriptionRequired: req.PrescriptionRequired,
	}

	if err := s.medicineRepo.Create(medicine); err != nil {
		return nil, err
	}

	return medicine, nil
}

func (s *medicineService) GetMedicineByID(id uint) (*models.Medicine, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	medicine, err := s.medicineRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if medicine == nil {
		return nil, ErrMedicineNotFound
	}

	return medicine, nil
}

func (s *medicineService) ListMedicines(filter models.MedicineFilter) ([]models.Medicine, int64, error) {
	return s.medicineRepo.FindAll(filter)
}

func (s *medicineService) UpdateMedicine(id uint, req models.MedicineUpdateRequest) (*models.Medicine, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}

	medicine, err := s.medicineRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if medicine == nil {
		return nil, ErrMedicineNotFound
	}

	if req.Name != nil {
		medicine.Name = *req.Name
	}
	if req.Description != nil {
		medicine.Description = *req.Description
	}
	if req.Price != nil {
		medicine.Price = *req.Price
	}
	if req.StockQuantity != nil {
		medicine.StockQuantity = *req.StockQuantity
	}
	if req.CategoryID != nil {
		category, err := s.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, ErrCategoryNotFound
		}
		medicine.CategoryID = *req.CategoryID
	}
	if req.SubCategoryID != nil {
		subCategory, err := s.subCategoryRepo.FindByID(*req.SubCategoryID)
		if err != nil {
			return nil, err
		}
		if subCategory == nil {
			return nil, ErrSubCategoryNotFound
		}

		categoryID := medicine.CategoryID
		if req.CategoryID != nil {
			categoryID = *req.CategoryID
		}

		if subCategory.CategoryID != categoryID {
			return nil, ErrInvalidSubCategory
		}

		medicine.SubCategoryID = req.SubCategoryID
	}
	if req.Manufacturer != nil {
		medicine.Manufacturer = *req.Manufacturer
	}
	if req.PrescriptionRequired != nil {
		medicine.PrescriptionRequired = *req.PrescriptionRequired
	}

	if err := s.medicineRepo.Update(medicine); err != nil {
		return nil, err
	}

	return medicine, nil
}

func (s *medicineService) DeleteMedicine(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}

	exists, err := s.medicineRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMedicineNotFound
	}

	return s.medicineRepo.Delete(id)
}

func (s *medicineService) CheckStock(medicineID uint, quantity int) (*models.Medicine, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	medicine, err := s.medicineRepo.FindByID(medicineID)
	if err != nil {
		return nil, err
	}
	if medicine == nil {
		return nil, ErrMedicineNotFound
	}

	if medicine.StockQuantity < quantity {
		return nil, fmt.Errorf("%w: requested %d, available %d",
			ErrInsufficientStock, quantity, medicine.StockQuantity)
	}

	return medicine, nil
}

func (s *medicineService) ReserveStock(medicineID uint, quantity int) error {
	medicine, err := s.CheckStock(medicineID, quantity)

	if err != nil {
		return err
	}

	newQuantity := medicine.StockQuantity - quantity
	return s.medicineRepo.UpdateStock(medicineID, newQuantity)
}

func (s *medicineService) ReleaseStock(medicineID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be positive")
	}

	medicine, err := s.medicineRepo.FindByID(medicineID)
	if err != nil {
		return err
	}
	if medicine == nil {
		return ErrMedicineNotFound
	}

	newQuantity := medicine.StockQuantity + quantity
	return s.medicineRepo.UpdateStock(medicineID, newQuantity)
}
