package userPort // secondary port

import (
	"context"

	"github.com/wittawat/go-hex/core/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindById(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	Find(ctx context.Context) ([]entities.User, error)
	UpdateOne(ctx context.Context, user *entities.User, id string) error
	DeleteOne(ctx context.Context, id string) error
}
