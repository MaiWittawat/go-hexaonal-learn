package services

import (
	"context"
	"log"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/validate"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	productPort "github.com/wittawat/go-hex/core/port/product"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"github.com/wittawat/go-hex/utils/errs"
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
	if err != nil { // case user not found
		log.Println("Error Create Product(findByEmail): ", err)
		return errs.ErrUserNotFound
	}

	if err := validator.IsValidProduct(product); err != nil {
		log.Println("Error Create Product(valid): ", err)
		return errs.ErrInValidProduct
	}
	product.CreatedBy = user.ID
	product.CreatedAt = time.Now()
	product.UpdatedAt = product.CreatedAt
	product.DeletedAt = nil
	if err := s.productRepo.Save(ctx, product); err != nil {
		log.Println("Error Create Product(db): ", err)
		return errs.ErrSaveProduct
	}
	return nil
}

func (s *ProductService) GetAll(ctx context.Context) ([]entities.Product, error) {
	products, err := s.productRepo.Find(ctx)
	if err != nil {
		log.Println("Error GetAll Product(find): ", err)
		return nil, errs.ErrProductNotFound
	}
	return products, nil
}

func (s *ProductService) GetById(ctx context.Context, id string) (*entities.Product, error) {
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error GetById Product(findById): ", err)
		return nil, errs.ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) EditOne(ctx context.Context, newProduct *entities.Product, id string, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Error EditOne Product(findByEmail): ", err)
		return errs.ErrUserNotFound
	}
	oldProduct, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error EditOne Product(findById): ", err)
		return errs.ErrProductNotFound
	}
	if user.ID != oldProduct.CreatedBy {
		log.Println("Error EditOne Product(userId): ", err)
		return errs.ErrForbidden
	}
	updateProduct, err := validator.EnsureUpdateProduct(oldProduct, newProduct)
	if err != nil {
		log.Println("Error EditOne Product(ensure): ", err)
		return errs.ErrInValidProduct
	}
	updateProduct.UpdatedAt = time.Now()
	if err := s.productRepo.UpdateOne(ctx, updateProduct, id); err != nil {
		log.Println("Error EditOne Product(db): ", err)
		return errs.ErrUpdateProduct
	}

	return nil
}

func (s *ProductService) DropOne(ctx context.Context, id string, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Error DropOne Product(findByEmail): ", err)
		return errs.ErrUserNotFound
	}
	product, err := s.productRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error DropOne Product(findById): ", err)
		return errs.ErrProductNotFound
	}
	if user.ID != product.CreatedBy {
		log.Println("Error DropOne Product(userId): ", err)
		return errs.ErrForbidden
	}
	if err := s.orderRepo.DeleteAllOrderByProduct(ctx, id); err != nil {
		log.Println("Error DropOne Product(orderDB): ", err)
		return errs.ErrDeleteAllOrderByProduct
	}
	if err := s.productRepo.DeleteOne(ctx, id); err != nil {
		log.Println("Error DropOne Product(productDB): ", err)
		return errs.ErrDeleteProduct
	}
	return nil
}
