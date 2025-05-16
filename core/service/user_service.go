package service

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/util/validate"
	"github.com/wittawat/go-hex/errs"

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
		log.Println("Error Save User(valid): ", err)
		return errs.ErrInValidUser
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	if err != nil {
		log.Println("Error Save User(hash): ", err)
		return errs.ErrHashPassword
	}

	user.Role = role
	user.Password = string(hash)
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.DeletedAt = nil

	if err := s.userRepo.Save(ctx, user); err != nil {
		log.Println("Error Save User(db): ", err)
		return errs.ErrSaveUser
	}
	return nil
}

func (s *UserService) UpdateOne(ctx context.Context, newUser *entities.User, id string, tokenEmail string) (string, jwt.MapClaims, error) {
	oldUser, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error UpdateOne User(findID): ", err)
		return "", nil, errs.ErrUserNotFound
	}

	updateUser, err := validator.IsValidUpdateUser(oldUser, newUser, tokenEmail)
	if err != nil {
		log.Println("Error UpdateOne User(valid): ", err)
		return "", nil, errs.ErrInValidUser
	}
	updateUser.UpdatedAt = time.Now()
	updateUser.Role = oldUser.Role
	var newToken string
	if updateUser.Email != tokenEmail {
		newToken, err = s.token.CreateToken(updateUser.Email)
		if err != nil {
			log.Println("Error UpdateOne User(Ctoken): ", err)
			return "", nil, errs.ErrCreatToken
		}
	}
	claims, err := s.token.VerifyToken(newToken)
	if err != nil {
		log.Println("Error UpdateOne User(Vtoken): ", err)
		return "", nil, errs.ErrVerifyToken
	}
	if err := s.userRepo.UpdateOne(ctx, updateUser, id); err != nil {
		log.Println("Error UpdateOne User(db): ", err)
		return "", nil, errs.ErrUpdateUser
	}
	return newToken, claims, nil
}

func (s *UserService) DeleteOne(ctx context.Context, id string, email string) error {
	existUser, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error DeleteOne User(findID): ", err)
		return errs.ErrUserNotFound
	}
	if existUser.Email != email {
		log.Println("Error DeleteOne User(email): ", err)
		return errs.ErrForbidden
	}

	if err := s.orderRepo.DeleteAllOrderByUser(ctx, id); err != nil {
		log.Println("Error DeleteOne User(orderDB): ", err)
		return errs.ErrDeleteAllOrderByUser
	}
	if err := s.productRepo.DeleteAll(ctx, id); err != nil {
		log.Println("Error DeleteOne User(productDB): ", err)
		return errs.ErrDeleteProduct
	}
	if err := s.userRepo.DeleteOne(ctx, id); err != nil {
		log.Println("Error DeleteOne User(userDB): ", err)
		return errs.ErrDeleteUser
	}
	return nil
}

func (s *UserService) Login(ctx context.Context, user *entities.User) (string, error) {
	u, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		log.Println("Error Login User(findByEmail): ", err)
		return "", errs.ErrLogin
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		log.Println("Error Login User(comparePass): ", err)
		return "", errs.ErrPassNotMatch
	}

	if user.Email != u.Email || user.Username != u.Username {
		log.Println("Error Login User(invalid): ", err)
		return "", errs.ErrInvalidInput
	}

	token, err := s.token.CreateToken(u.Email)
	if err != nil {
		log.Println("Error Login User(token): ", err)
		return "", errs.ErrCreatToken
	}

	return token, nil
}

func (s *UserService) Find(ctx context.Context) ([]entities.User, error) {
	users, err := s.userRepo.Find(ctx)
	if err != nil {
		log.Println("Error Find User(db): ", err)
		return nil, errs.ErrUserNotFound
	}
	return users, nil
}

func (s *UserService) FindById(ctx context.Context, id string) (*entities.User, error) {
	user, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error FindById User(db): ", err)
		return nil, errs.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Error FindByEmail User(db): ", err)
		return nil, errs.ErrUserNotFound
	}
	return user, nil
}
