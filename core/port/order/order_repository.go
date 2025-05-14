package orderPort

import (
	"context"

	"github.com/wittawat/go-hex/core/entities"
)

// outbound
type OrderRepository interface {
	Save(ctx context.Context, order *entities.Order) error
	FindByUserId(ctx context.Context, userId string) ([]entities.Product, error)
	FindByUserEmail(ctx context.Context, email string) (*entities.Order, error)
	FindById(ctx context.Context, id string) (*entities.Order, error)
	UpdateOne(ctx context.Context, order *entities.Order, id string) error
	DeleteOne(ctx context.Context, id string) error
	DeleteAllByUser(ctx context.Context, userId string) error
	DeleteAllByProduct(ctx context.Context, productId string) error
}
