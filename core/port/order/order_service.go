package orderPort

import "github.com/wittawat/go-hex/core/entities"

// inbound
type OrderService interface {
	Create(order *entities.Order) error
	GetByUser(userId int) ([]entities.Product, error)
	Update(order *entities.Order, id int) error
	Delete(id int) error
}
