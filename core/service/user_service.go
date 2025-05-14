package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/util/validate"

	authPort "github.com/wittawat/go-hex/core/port/auth"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	productPort "github.com/wittawat/go-hex/core/port/product"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo    userPort.UserRepository //user repository
	productRepo productPort.ProductRepository
	orderRepo   orderPort.OrderRepository
	token       authPort.JwtAuthNService
}

func NewUserService(useRepo userPort.UserRepository, productRepo productPort.ProductRepository, orderRepo orderPort.OrderRepository, token authPort.JwtAuthNService) userPort.UserService {
	return &UserService{userRepo: useRepo, productRepo: productRepo, orderRepo: orderRepo, token: token}
}

func (s *UserService) Save(ctx context.Context, user *entities.User, role string) error {

	if err := validator.IsValidUser(user); err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		return err
	}

	user.Role = role
	user.Password = string(hash)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.DeletedAt = nil

	if err := s.userRepo.Save(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateOne(ctx context.Context, newUser *entities.User, id string, tokenEmail string) (string, jwt.MapClaims, error) {
	oldUser, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return "", nil, err
	}

	updateUser, err := validator.IsValidUpdateUser(oldUser, newUser, tokenEmail)
	if err != nil {
		return "", nil, err
	}
	updateUser.UpdatedAt = time.Now()
	updateUser.Role = "user"
	var newToken string
	if updateUser.Email != tokenEmail {
		newToken, err = s.token.CreateToken(updateUser.Email)
		if err != nil {
			return "", nil, err
		}
	}
	claims, err := s.token.VerifyToken(newToken)
	if err != nil {
		return "", nil, err
	}
	if err := s.userRepo.UpdateOne(ctx, updateUser, id); err != nil {
		return "", nil, err
	}
	return newToken, claims, nil
}

func (s *UserService) DeleteOne(ctx context.Context, id string, email string) error {
	existUser, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if existUser.Email != email {
		return errors.New("fail to delete user wrong email")
	}
	if err := s.orderRepo.DeleteAllByUser(ctx, id); err != nil {
		return err
	}
	if err := s.productRepo.DeleteAll(ctx, id); err != nil {
		return err
	}
	if err := s.userRepo.DeleteOne(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, user *entities.User) (string, error) {
	u, err := s.userRepo.FindByEmail(ctx, user.Email)
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

func (s *UserService) Find(ctx context.Context) ([]entities.User, error) {
	users, err := s.userRepo.Find(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
