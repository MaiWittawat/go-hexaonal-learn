package service

import (
	"errors"

	"github.com/wittawat/go-hex/core/entities"
	authPort "github.com/wittawat/go-hex/core/port/auth"
	userPort "github.com/wittawat/go-hex/core/port/user"
)

type UserService struct {
	repo  userPort.UserRepository //user repository
	token authPort.JwtAuthNService
}

func NewUserService(repo userPort.UserRepository, token authPort.JwtAuthNService) userPort.UserService {
	return &UserService{repo: repo, token: token}
}

func (s *UserService) Save(user *entities.User) error {
	if len(user.Password) < 4 {
		return errors.New("invalid password")
	}

	if err := s.repo.Save(user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) FindById(id int) (*entities.User, error) {
	user, err := s.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByEmail(email string) (*entities.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Find() ([]entities.User, error) {
	users, err := s.repo.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) UpdateOne(user *entities.User, id int) error {
	if err := s.repo.UpdateOne(user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteOne(id int) error {
	if err := s.repo.DeleteOne(id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(user *entities.User) (string, error) {
	u, err := s.repo.FindByEmail(user.Email)
	if err == nil && u.Email == user.Email && u.Password == user.Password {
		token, err := s.token.CreateToken(u.Email)
		if err != nil {
			return "", nil
		}
		return token, nil
	}
	return "", nil
}
