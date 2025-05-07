package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"default:\"guest\""`
	Role     string `gorm:"default:\"user\""`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
}
