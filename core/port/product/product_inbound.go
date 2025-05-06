package port

import "github.com/wittawat/go-hex/core/entities"

type ProductInbound interface {
	Save(product *entities.Product) error
	FindById(id int) (*entities.Product, error)
	Find() ([]entities.Product, error)
	UpdateOne(product *entities.Product, id int) error
	DeleteOne(id int) error
}
