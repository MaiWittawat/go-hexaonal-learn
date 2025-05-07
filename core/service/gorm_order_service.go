package service

import (
	"github.com/wittawat/go-hex/core/entities"
	port "github.com/wittawat/go-hex/core/port/order"
)

type GormOrderService struct {
	repo port.OrderRepository
}

func NewGormOrderService(repo port.OrderRepository) port.OrderService {
	return &GormOrderService{repo: repo}
}

func (s *GormOrderService) Create(order *entities.Order) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}
	return nil
}

func (s *GormOrderService) GetByUser(userId int) ([]entities.Product, error) {
	products, err := s.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *GormOrderService) Update(order *entities.Order, id int) error {
	if err := s.repo.UpdateOne(order, id); err != nil {
		return err
	}
	return nil
}

func (s *GormOrderService) Delete(id int) error {
	if err := s.repo.DeleteOne(id); err != nil {
		return err
	}
	return nil
}
