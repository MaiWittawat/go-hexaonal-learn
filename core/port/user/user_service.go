package userPort // primary port

import (
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/entities/request"
)

type UserService interface {
	Save(user *entities.User) error
	FindById(id int) (*entities.User, error)
	Find() ([]entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	UpdateOne(user *entities.User, id int, email string) error
	DeleteOne(id int, email string) error
	Login(user *request.UserRequest) (string, error)
}
