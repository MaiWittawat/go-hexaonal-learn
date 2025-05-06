package port // secondary port

import "github.com/wittawat/go-hex/core/entities"

type UserOutbound interface {
	Save(user *entities.User) error
	FindById(id int) (*entities.User, error)
	Find() ([]entities.User, error)
	UpdateOne(user *entities.User, id int) error
	DeleteOne(id int) error
}
