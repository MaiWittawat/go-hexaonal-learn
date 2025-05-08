package userPort // secondary port

import "github.com/wittawat/go-hex/core/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindById(id int) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	Find() ([]entities.User, error)
	UpdateOne(user *entities.User, id int) error
	DeleteOne(id int) error
}
