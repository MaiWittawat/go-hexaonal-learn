package entities

import "time"

type Order struct {
	ID        string
	UserID    string
	ProductID string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
