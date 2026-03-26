package models

import "gorm.io/gorm"

// cart
// item_id -- user_id -- medicine_id -- quantity
//   1     --    1    --     5       --    3
// 2 -- 1 -- 8  -- 2
// 3 -- 1 -- 10 -- 6

// cart itsmm
//  -- medicine_id - quantity
//

// medicines
// medicine_id
// price_per_unit
// ...

// users
// user_id
// ...

type CartItem struct {
	gorm.Model
	MedicineID   int  `json:"medicine_id"`
	Quantity     int  `json:"quantity"`
	LineTotal    int  `json:"line_total"`
	PricePerUnit int  `json:"price_per_unit"`
	CartID       uint `json:"cart_id"`
}

type CartItemCreateRequest struct {
	PricePerUnit int `json:"price_per_unit" binding:"required"`
	MedicineID   int `json:"medicine_id" binding:"required"`
	Quantity     int `json:"quantity" binding:"required"`
}

type CartItemUpdateRequest struct {
	Quantity   *int `json:"quantity"`
}

// ---------
// 1 - парацетомол (5 id)
// 2
