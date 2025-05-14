package userPort // primary port

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wittawat/go-hex/core/entities"
)

type UserService interface {
	Save(ctx context.Context, user *entities.User, role string) error
	FindById(ctx context.Context, id string) (*entities.User, error)
	Find(ctx context.Context) ([]entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateOne(ctx context.Context, user *entities.User, id string, email string) (string, jwt.MapClaims, error)
	DeleteOne(ctx context.Context, id string, email string) error
	Login(ctx context.Context, user *entities.User) (string, error)
}
