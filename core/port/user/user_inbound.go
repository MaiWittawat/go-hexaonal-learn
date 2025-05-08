package port // primary port

import "github.com/wittawat/go-hex/core/entities"

type UserInbound interface {
	Save(user *entities.User) error
	FindById(id int) (*entities.User, error)
	Find() ([]entities.User, error)
	FindByEmail(email string) (*entities.User, error)
	UpdateOne(user *entities.User, id int) error
	DeleteOne(id int) error
	Login(user *entities.User) (string, error)
}
