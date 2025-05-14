package service

import (
	"context"
	"errors"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/util/validate"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	productPort "github.com/wittawat/go-hex/core/port/product"
	userPort "github.com/wittawat/go-hex/core/port/user"
)

var (
	ErrProductNotFound = errors.New("Product not found")
	ErrUserNotFound    = errors.New("User not found")
	ErrForbidden       = errors.New("No permission")
)

type ProductService struct {
	userRepo    userPort.UserRepository
	productRepo productPort.ProductRepository
	orderRepo   orderPort.OrderRepository
}

func NewProductService(userRepo userPort.UserRepository, productRepo productPort.ProductRepository, orderRepo orderPort.OrderRepository) productPort.ProductService {
	return &ProductService{userRepo: userRepo, productRepo: productRepo, orderRepo: orderRepo}
}

func (s *ProductService) Create(ctx context.Context, product *entities.Product, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil { // case use not found
		return ErrUserNotFound
	}

	if err := validator.IsValidProduct(product); err != nil {
		return err
	}
	product.CreatedBy = user.ID
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt
	product.DeletedAt = nil
	if err := s.productRepo.Save(ctx, product); err != nil {
		return err
	}
	return nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]entities.Product, error) {
	products, err := s.productRepo.Find(ctx)
	if err != nil {
		return nil, ErrProductNotFound
	}
	return products, nil
}

func (s *ProductService) GetById(ctx context.Context, id string) (*entities.Product, error) {
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return product, nil
}

func (s *ProductService) EditOne(ctx context.Context, newProduct *entities.Product, id string, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	oldProduct, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return ErrProductNotFound
	}
	if user.ID != oldProduct.CreatedBy {
		return ErrForbidden
	}
	updateProduct, err := validator.EnsureUpdateProduct(oldProduct, newProduct)
	if err != nil {
		return err
	}
	updateProduct.UpdatedAt = time.Now()
	if err := s.productRepo.UpdateOne(ctx, updateProduct, id); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) DropOne(ctx context.Context, id string, email string) error {
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		return ErrProductNotFound
	}
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return ErrUserNotFound
	}
	if user.ID != product.CreatedBy {
		return ErrForbidden
	}
	if err := s.orderRepo.DeleteAllByProduct(ctx, id); err != nil {
		return err
	}
	if err := s.productRepo.DeleteOne(ctx, id); err != nil {
		return err
	}
	return nil
}
