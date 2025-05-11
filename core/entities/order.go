package entities

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId    uint `gorm:"user_id"`
	ProductId uint `gorm:"product_id"`
}
