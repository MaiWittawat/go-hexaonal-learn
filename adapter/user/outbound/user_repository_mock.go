package userAdapter

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wittawat/go-hex/core/entities"
)

type userRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *userRepositoryMock {
	return &userRepositoryMock{}
}

func (m *userRepositoryMock) Save(ctx context.Context, user *entities.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *userRepositoryMock) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *userRepositoryMock) FindById(ctx context.Context, id string) (*entities.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), nil
}

func (m *userRepositoryMock) Find(ctx context.Context) ([]entities.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entities.User), nil
}

func (m *userRepositoryMock) UpdateOne(ctx context.Context, user *entities.User, id string) error {
	args := m.Called(ctx, user, id)
	return args.Error(0)
}

func (m *userRepositoryMock) DeleteOne(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}
