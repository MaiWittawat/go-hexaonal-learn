package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title  string `json:"title"`
	Price  uint   `json:"price"`
	Detail string `json:"detail"`
}
