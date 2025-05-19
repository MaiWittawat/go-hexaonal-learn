package productAdapter

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wittawat/go-hex/core/entities"
)

type productRepositoryMock struct {
	mock.Mock
}

func NewProductRepositoryMock() *productRepositoryMock {
	return &productRepositoryMock{}
}

func (m *productRepositoryMock) Save(ctx context.Context, product *entities.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *productRepositoryMock) FindById(ctx context.Context, id string) (*entities.Product, error) {
	args := m.Called(ctx, id)
	if temp := args.Get(0); temp != nil {
		return temp.(*entities.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *productRepositoryMock) Find(ctx context.Context) ([]entities.Product, error) {
	args := m.Called(ctx)
	if temp := args.Get(0); temp != nil {
		return temp.([]entities.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *productRepositoryMock) UpdateOne(ctx context.Context, product *entities.Product, id string) error {
	args := m.Called(ctx, product, id)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (m *productRepositoryMock) DeleteOne(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *productRepositoryMock) DeleteAll(ctx context.Context, userId string) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
