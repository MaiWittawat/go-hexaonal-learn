package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	validator "github.com/wittawat/go-hex/core/internal/util/validate"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	userPort "github.com/wittawat/go-hex/core/port/user"
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
		return err
	}
	order.UserID = user.ID
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	fmt.Println("order in service: ", order)
	if err := s.orderRepo.Save(ctx, order); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetByUser(ctx context.Context, userId string) ([]entities.Product, error) {
	products, err := s.orderRepo.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *OrderService) EditOne(ctx context.Context, newOrder *entities.Order, id string, email string) error {
	oldOrder, err := s.orderRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.ID != oldOrder.UserID {
		return errors.New("you can update only you order")
	}
	order, err := validator.EnsureUpdateOrder(oldOrder, newOrder)
	if err != nil {
		return err
	}

	if err := s.orderRepo.UpdateOne(ctx, order, id); err != nil {
		return err
	}
	return nil
}

func (s *OrderService) DropOne(ctx context.Context, id string, email string) error {
	oldOrder, err := s.orderRepo.FindById(ctx, id)
	if err != nil {
		return err
	}
	o, err := s.orderRepo.FindByUserEmail(ctx, email)
	if err != nil {
		return err
	}
	if o.UserID != oldOrder.UserID {
		return errors.New("you can delete only you order")
	}
	if err := s.orderRepo.DeleteOne(ctx, id); err != nil {
		return err
	}
	return nil
}
