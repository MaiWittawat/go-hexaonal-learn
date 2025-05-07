package port

import "github.com/wittawat/go-hex/core/entities"

// outbound
type OrderRepository interface {
	Save(order *entities.Order) error
	FindByUserId(userId int) ([]entities.Product, error)
	FindById(id int) (*entities.Order, error)
	UpdateOne(order *entities.Order, id int) error
	DeleteOne(id int) error
}
