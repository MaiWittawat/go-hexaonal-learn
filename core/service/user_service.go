package service

import (
	"errors"

	"github.com/wittawat/go-hex/core/entities"
	port "github.com/wittawat/go-hex/core/port/user"
)

type UserService struct {
	ob port.UserOutbound //user repository
}

func NewUserService(ob port.UserOutbound) port.UserInbound {
	return &UserService{ob: ob}
}

func (s *UserService) Save(user *entities.User) error {
	if len(user.Password) < 4 {
		return errors.New("invalid password")
	}

	if err := s.ob.Save(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindById(id int) (*entities.User, error) {
	user, err := s.ob.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Find() ([]entities.User, error) {
	users, err := s.ob.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateOne(user *entities.User, id int) error {
	if err := s.ob.UpdateOne(user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteOne(id int) error {
	if err := s.ob.DeleteOne(id); err != nil {
		return err
	}
	return nil
}
