package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"username"`
	Role     string `gorm:"default:\"user\""`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"password"`
}
