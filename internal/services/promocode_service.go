package services

import (
	"errors"
	"pharmacy/internal/models"
	"pharmacy/internal/repository"

	"gorm.io/gorm"
)

type PromocodeService struct {
	PromoRepo repository.PromocodeRepository
}

func NewPromocodeService(
	promoRepo repository.PromocodeRepository,
) PromocodeService {
	return PromocodeService{
		PromoRepo: promoRepo,
	}
}

func (s *PromocodeService) Create(req models.PromocodeCreateRequest) (*models.Promocode, error) {
	existing, err := s.PromoRepo.GetByCode(req.Code)
	if err == nil && existing != nil {
		return nil, errors.New("такой код уже есть")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if req.ValidFrom != nil && req.ValidTo != nil && req.ValidTo.Before(*req.ValidFrom) {
		return nil, errors.New("valid_to не может быть раньше valid_from")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	promocode := &models.Promocode{
		Code:          req.Code,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		IsActive:      isActive,
	}

	if req.Description != nil {
		promocode.Description = *req.Description
	}
	if req.ValidFrom != nil {
		promocode.ValidFrom = *req.ValidFrom
	}
	if req.ValidTo != nil {
		promocode.ValidTo = *req.ValidTo
	}
	if req.MaxUses != nil {
		promocode.MaxUses = *req.MaxUses
	}
	if req.MaxUsesPerUser != nil {
		promocode.MaxUsesPerUser = *req.MaxUsesPerUser
	}

	if err := s.PromoRepo.Create(promocode); err != nil {
		return nil, err
	}

	return promocode, nil
}

func (s *PromocodeService) GetAll() ([]models.Promocode, error) {
	promocodes, err := s.PromoRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return promocodes, nil
}

func (s *PromocodeService) GetByCode(code string) (*models.Promocode, error) {
	promocode, err := s.PromoRepo.GetByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("промокод не найден")
		}
		return nil, err
	}
	if promocode == nil {
		return nil, errors.New("промокод не найден")
	}
	return promocode, nil
}

func (s *PromocodeService) Update(id uint, req *models.PromocodeUpdateRequest) (*models.Promocode, error) {
	if req == nil {
		return nil, errors.New("пустой запрос на обновление")
	}

	existing, err := s.PromoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("промокод не найден")
		}
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("промокод не найден")
	}

	if req.Code != nil && *req.Code != existing.Code {
		other, err := s.PromoRepo.GetByCode(*req.Code)
		if err == nil && other != nil && other.ID != existing.ID {
			return nil, errors.New("такой код уже есть")
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	if err := s.PromoRepo.Update(id, req); err != nil {
		return nil, err
	}

	return s.PromoRepo.GetByID(id)
}

func (s *PromocodeService) Delete(id uint) error {
	_, err := s.PromoRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("промокод не найден")
		}
		return err
	}

	return s.PromoRepo.Delete(id)
}
