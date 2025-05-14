package entities

import "time"

type Product struct {
	ID        string
	Title     string
	Price     int32
	Detail    string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
