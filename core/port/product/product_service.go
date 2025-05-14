package productPort

import (
	"context"

	"github.com/wittawat/go-hex/core/entities"
)

type ProductService interface {
	Create(ctx context.Context, product *entities.Product, emil string) error
	GetById(ctx context.Context, id string) (*entities.Product, error)
	GetAll(ctx context.Context) ([]entities.Product, error)
	EditOne(ctx context.Context, product *entities.Product, id string, email string) error
	DropOne(ctx context.Context, id string, email string) error
}
