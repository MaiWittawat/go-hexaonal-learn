package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	orderAdapter "github.com/wittawat/go-hex/adapter/order/outbound"
	userAdapter "github.com/wittawat/go-hex/adapter/user/outbound"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/services"
	"github.com/wittawat/go-hex/utils/errs"
)

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	orderService := services.NewOrderService(orderRepoMock, userRepoMock)

	t.Run("error_user_findByEmail", func(t *testing.T) {
		user := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com"}
		order := &entities.Order{UserID: "1", ProductID: "1"}

		userRepoMock.On("FindByEmail", ctx, user.Email).Return(nil, errs.ErrUserNotFound)

		err := orderService.Create(ctx, order, user.Email)

		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_save", func(t *testing.T) {
		user := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com"}
		order := &entities.Order{UserID: "1", ProductID: "1"}

		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)
		orderRepoMock.On("Save", ctx, order).Return(errs.ErrSaveOrder)

		err := orderService.Create(ctx, order, user.Email)

		assert.EqualError(t, err, errs.ErrSaveOrder.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		user := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com"}
		order := &entities.Order{UserID: "1", ProductID: "1"}

		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)
		orderRepoMock.On("Save", ctx, order).Return(nil)

		err := orderService.Create(ctx, order, user.Email)

		assert.Nil(t, err)

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})
}

func TestGetOrderByUser(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	orderService := services.NewOrderService(orderRepoMock, userRepoMock)

	t.Run("error_order_findByUserId", func(t *testing.T) {
		user := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com"}
		orderRepoMock.On("FindByUserId", ctx, user.ID).Return(nil, errs.ErrOrderNotFound)

		_, err := orderService.GetByUser(ctx, user.ID)
		assert.EqualError(t, err, errs.ErrOrderNotFound.Error())

		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		expectedUser := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com"}
		expectedProducts := []entities.Product{
			{ID: "1", Title: "product_title1", Price: 100, Detail: "product_detail1"},
			{ID: "2", Title: "product_title2", Price: 100, Detail: "product_detail2"},
			{ID: "3", Title: "product_title3", Price: 100, Detail: "product_detail3"},
		}

		orderRepoMock.On("FindByUserId", ctx, expectedUser.ID).Return(expectedProducts, nil)

		products, _ := orderService.GetByUser(ctx, expectedUser.ID)
		assert.Equal(t, expectedProducts, products)
		orderRepoMock.ExpectedCalls = nil
	})
}

func TestEditOneOrder(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	orderService := services.NewOrderService(orderRepoMock, userRepoMock)

	t.Run("error_order_findById", func(t *testing.T) {
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		orderRepoMock.On("FindById", ctx, order.ID).Return(nil, errs.ErrOrderNotFound)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrOrderNotFound.Error())

		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_findByEmail", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(order, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(nil, errs.ErrUserNotFound)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_idNotMatch", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		expectedOrder := &entities.Order{ID: "2", UserID: "2", ProductID: "2"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(expectedOrder, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrForbidden.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_updateOne", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		oldOrder := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		newOrder := &entities.Order{ID: "1", UserID: "1", ProductID: "2"}

		orderRepoMock.On("FindById", ctx, oldOrder.ID).Return(oldOrder, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)
		orderRepoMock.On("UpdateOne", ctx, mock.AnythingOfType("*entities.Order"), oldOrder.ID).Return(errs.ErrUpdateOrder)

		err := orderService.EditOne(ctx, newOrder, oldOrder.ID, user.Email)
		assert.EqualError(t, err, errs.ErrUpdateOrder.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		newOrder := &entities.Order{ID: "1", UserID: "1", ProductID: "2"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(order, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)
		orderRepoMock.On("UpdateOne", ctx, mock.AnythingOfType("*entities.Order"), order.ID).Return(nil)

		err := orderService.EditOne(ctx, newOrder, order.ID, "mai@example.com")

		assert.Nil(t, err)
		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})
}

func TestDropOneOrder(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	orderService := services.NewOrderService(orderRepoMock, userRepoMock)

	t.Run("error_order_findById", func(t *testing.T) {
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		orderRepoMock.On("FindById", ctx, order.ID).Return(nil, errs.ErrOrderNotFound)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrOrderNotFound.Error())

		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_findByEmail", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(order, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(nil, errs.ErrUserNotFound)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_idNotMatch", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}
		expectedOrder := &entities.Order{ID: "2", UserID: "2", ProductID: "2"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(expectedOrder, nil)
		userRepoMock.On("FindByEmail", ctx, user.Email).Return(user, nil)

		err := orderService.EditOne(ctx, order, order.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrForbidden.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_order_deleteOne", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(order, nil)
		orderRepoMock.On("FindByUserEmail", ctx, user.Email).Return(order, nil)
		orderRepoMock.On("DeleteOne", ctx, order.ID).Return(errs.ErrDeleteOrder)

		err := orderService.DropOne(ctx, "1", "mai@example.com")
		assert.EqualError(t, err, errs.ErrDeleteOrder.Error())

		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com"}
		order := &entities.Order{ID: "1", UserID: "1", ProductID: "1"}

		orderRepoMock.On("FindById", ctx, order.ID).Return(order, nil)
		orderRepoMock.On("FindByUserEmail", ctx, user.Email).Return(order, nil)
		orderRepoMock.On("DeleteOne", ctx, order.ID).Return(nil)

		err := orderService.DropOne(ctx, "1", "mai@example.com")
		assert.Nil(t, err)

		orderRepoMock.ExpectedCalls = nil
	})
}
