package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId    uint `json:"user_id"`
	ProductId uint `json:"product_id"`
}
