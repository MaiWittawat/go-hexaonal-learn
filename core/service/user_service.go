package service

import (
	"errors"

	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/entities/request"
	authPort "github.com/wittawat/go-hex/core/port/auth"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"golang.org/x/crypto/bcrypt"
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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		return err
	}

	user.Password = string(hash)
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

func (s *UserService) UpdateOne(user *entities.User, id int, email string) error {
	existUser, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	if existUser.Email != email {
		return errors.New("fail to update user wrong email")
	}

	if user.Username == "" {
		user.Username = existUser.Username
	}
	if user.Email == "" {
		user.Email = existUser.Email
	}
	if user.Password == "" {
		user.Password = existUser.Password
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
		if err != nil {
			return err
		}
		user.Password = string(hash)
	}

	if err := s.repo.UpdateOne(user, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteOne(id int, email string) error {
	existUser, err := s.repo.FindById(id)
	if err != nil {
		return err
	}
	if existUser.Email != email {
		return errors.New("fail to delete user wrong email")
	}

	if err := s.repo.DeleteOne(id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(user *request.UserRequest) (string, error) {
	u, err := s.repo.FindByEmail(user.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return "", errors.New("password not match")
	}

	if user.Email != u.Email || user.Username != u.Username {
		return "", errors.New("invalid fail to login")
	}

	token, err := s.token.CreateToken(u.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
