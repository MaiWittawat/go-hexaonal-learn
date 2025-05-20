package services

import (
	"context"
	"log"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/validate"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	userPort "github.com/wittawat/go-hex/core/port/user"
	"github.com/wittawat/go-hex/utils/errs"
)

type OrderService struct {
	orderRepo orderPort.OrderRepository
	userRepo  userPort.UserRepository
}

func NewOrderService(orderRepo orderPort.OrderRepository, userRepo userPort.UserRepository) orderPort.OrderService {
	return &OrderService{orderRepo: orderRepo, userRepo: userRepo}
}

func (s *OrderService) Create(ctx context.Context, order *entities.Order, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Error Create Order(findByEmail): ", err)
		return errs.ErrUserNotFound
	}
	order.UserID = user.ID
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	if err := s.orderRepo.Save(ctx, order); err != nil {
		log.Println("Error Create Order(db): ", err)
		return errs.ErrSaveOrder
	}
	return nil
}

func (s *OrderService) GetByUser(ctx context.Context, userId string) ([]entities.Product, error) {
	products, err := s.orderRepo.FindByUserId(ctx, userId)
	if err != nil {
		log.Println("Error GetByUser Order(findByUserId): ", err)
		return nil, errs.ErrOrderNotFound
	}
	return products, nil
}

func (s *OrderService) EditOne(ctx context.Context, newOrder *entities.Order, id string, email string) error {
	oldOrder, err := s.orderRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error EditOne Order(findById): ", err)
		return errs.ErrOrderNotFound
	}
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Println("Error EditOne Order(findByEmail): ", err)
		return errs.ErrUserNotFound
	}

	if user.ID != oldOrder.UserID {
		log.Println("Error EditOne Order(userId): ", err)
		return errs.ErrForbidden
	}
	updateOrder := validator.EnsureUpdateOrder(oldOrder, newOrder)

	if err := s.orderRepo.UpdateOne(ctx, updateOrder, id); err != nil {
		log.Println("Error EditOne Order(db): ", err)
		return errs.ErrUpdateOrder
	}
	return nil
}

func (s *OrderService) DropOne(ctx context.Context, id string, email string) error {
	oldOrder, err := s.orderRepo.FindById(ctx, id)
	if err != nil {
		log.Println("Error DropOne Order(findById): ", err)
		return errs.ErrOrderNotFound
	}
	o, err := s.orderRepo.FindByUserEmail(ctx, email)
	if err != nil {
		log.Println("Error DropOne Order(findByEmail): ", err)
		return errs.ErrOrderNotFound
	}
	if o.UserID != oldOrder.UserID {
		log.Println("Error DropOne Order(userId): ", err)
		return errs.ErrForbidden
	}
	if err := s.orderRepo.DeleteOne(ctx, id); err != nil {
		log.Println("Error DropOne Order(db): ", err)
		return errs.ErrDeleteOrder
	}
	return nil
}
