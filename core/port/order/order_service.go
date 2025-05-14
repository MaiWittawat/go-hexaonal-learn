package orderPort

import (
	"context"

	"github.com/wittawat/go-hex/core/entities"
)

// inbound
type OrderService interface {
	Create(ctx context.Context, order *entities.Order, email string) error
	GetByUser(ctx context.Context, userId string) ([]entities.Product, error)
	EditOne(ctx context.Context, order *entities.Order, id string, email string) error
	DropOne(ctx context.Context, id string, email string) error
}
