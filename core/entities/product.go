package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title  string `gorm:"title"`
	Price  uint   `gorm:"price"`
	Detail string `gorm:"detail"`
}
