package repository

import (
	"pharmacy/internal/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	Create(card *models.Cart) error
	GetByUserID(userID uint64) (*models.Cart, error)

	AddItem(cardItem *models.CartItem) error

	UpdateItem(item *models.CartItem) error
	DeleteItem(itemID uint64) error
	ClearCart(userID uint64) error
}

type gormCartRepository struct {
	db *gorm.DB
}

// func NewCardRepository(db *gorm.DB) CartRepository {
// 	return &gormCartRepository{db: db}
// }

func (r gormCartRepository) ClearCart(userID uint64) error {
	cart, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	return r.db.Model(cart).Association("CartItems").Clear()

}

func (r gormCartRepository) DeleteItem(itemID uint64) error {
	return r.db.Delete(&models.CartItem{}, itemID).Error
}

func (r gormCartRepository) UpdateItem(item *models.CartItem) error {
	return r.db.Model(&models.CartItem{}).Where("id = ?", item.ID).Select("*").Updates(item).Error
}

func (r gormCartRepository) Create(card *models.Cart) error {
	return r.db.Create(card).Error
}

func (r gormCartRepository) GetByUserID(userID uint64) (*models.Cart, error) {
	var cart models.Cart
	if err := r.db.First(&cart, userID).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r gormCartRepository) AddItem(userID uint, cardItem *models.CartItem) error {
	// получить Cart
	var card *models.Cart
	if err := r.db.First(&card, cardItem.CartID).Error; err != nil {
		card = &models.Cart{
			UserID: userID,
		}
		r.Create(card)
	}
	// если Cart нет - создать

	// Получить созданный Cart

	item := models.CartItem{
		MedicineID: cardItem.MedicineID,
		Quantity:   cardItem.Quantity,
		CartID:     card.ID,
		// <- вроде как не нужен, если мои работаем через Association (тк Cart уже указано)
		// учесть следующие поля (Quantity, ..., LineTotal)
	}

	err := r.db.Model(&card).Association("CartItems").Append(&item)

	return err
}

// func (r gormCartRepository)
