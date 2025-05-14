package productPort

import (
	"context"

	"github.com/wittawat/go-hex/core/entities"
)

type ProductRepository interface {
	Save(ctx context.Context, product *entities.Product) error
	FindById(ctx context.Context, id string) (*entities.Product, error)
	Find(ctx context.Context) ([]entities.Product, error)
	UpdateOne(ctx context.Context, product *entities.Product, id string) error
	DeleteOne(ctx context.Context, id string) error
	DeleteAll(ctx context.Context, userId string) error
}
