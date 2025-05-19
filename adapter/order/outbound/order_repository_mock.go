package orderAdapter

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wittawat/go-hex/core/entities"
)

type orderRepositoryMock struct {
	mock.Mock
}

func NewOrderRepositoryMock() *orderRepositoryMock {
	return &orderRepositoryMock{}
}

func (m *orderRepositoryMock) Save(ctx context.Context, order *entities.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *orderRepositoryMock) FindByUserId(ctx context.Context, userId string) ([]entities.Product, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.Product), nil
}

func (m *orderRepositoryMock) FindByUserEmail(ctx context.Context, email string) (*entities.Order, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *orderRepositoryMock) FindById(ctx context.Context, id string) (*entities.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Order), nil
}

func (m *orderRepositoryMock) UpdateOne(ctx context.Context, order *entities.Order, id string) error {
	args := m.Called(ctx, order, id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *orderRepositoryMock) DeleteOne(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *orderRepositoryMock) DeleteAllOrderByUser(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *orderRepositoryMock) DeleteAllOrderByProduct(ctx context.Context, productId string) error {
	args := m.Called(ctx, productId)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
