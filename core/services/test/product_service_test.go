package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	orderAdapter "github.com/wittawat/go-hex/adapter/order/outbound"
	productAdapter "github.com/wittawat/go-hex/adapter/product/outbound"
	userAdapter "github.com/wittawat/go-hex/adapter/user/outbound"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/services"
	"github.com/wittawat/go-hex/utils/errs"
)

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	productService := services.NewProductService(userRepoMock, productRepoMock, orderRepoMock)

	t.Run("success", func(t *testing.T) {
		user := &entities.User{ID: "123", Email: "test@example.com"}
		product := &entities.Product{Title: "test", Price: 100, Detail: "t te tes test detail"}

		userRepoMock.On("FindByEmail", ctx, "test@example.com").Return(user, nil)
		productRepoMock.On("Save", ctx, mock.MatchedBy(func(p *entities.Product) bool {
			return p.Title == "test" && p.CreatedBy == "123"
		})).Return(nil)

		err := productService.Create(ctx, product, "test@example.com")
		require.NoError(t, err)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_notfound", func(t *testing.T) {
		product := &entities.Product{Title: "test", Price: 100, Detail: "desc"}
		userRepoMock.On("FindByEmail", mock.Anything, "unknown@example.com").Return(nil, errs.ErrUserNotFound)

		err := productService.Create(ctx, product, "unknown@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_product_invalid", func(t *testing.T) {
		user := &entities.User{ID: "123", Email: "test@example.com"}
		product := &entities.Product{Title: "", Price: -1, Detail: ""}

		userRepoMock.On("FindByEmail", ctx, "test@example.com").Return(user, nil)

		err := productService.Create(ctx, product, user.Email)
		assert.EqualError(t, err, errs.ErrInValidProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_product_save", func(t *testing.T) {
		user := &entities.User{ID: "123", Email: "test@example.com"}
		product := &entities.Product{Title: "apple", Price: 20, Detail: "Just a normal sweet apple"}

		userRepoMock.On("FindByEmail", ctx, "test@example.com").Return(user, nil)
		productRepoMock.On("Save", mock.Anything, product).Return(errs.ErrSaveProduct)
		err := productService.Create(ctx, product, user.Email)
		assert.EqualError(t, err, errs.ErrSaveProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

}

func TestGetByIdProduct(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	productService := services.NewProductService(userRepoMock, productRepoMock, orderRepoMock)

	t.Run("success", func(t *testing.T) {
		expectedProduct := &entities.Product{Title: "Product A", Price: 10, Detail: "Product A details"}
		productRepoMock.On("FindById", mock.Anything, "123").Return(expectedProduct, nil)

		product, _ := productService.GetById(ctx, "123")
		assert.Equal(t, product, expectedProduct)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_product_findById", func(t *testing.T) {
		productRepoMock.On("FindById", mock.Anything, "2").Return(nil, errs.ErrProductNotFound)
		_, err := productService.GetById(ctx, "2")
		assert.EqualError(t, err, errs.ErrProductNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})
}

func TestGetAllProduct(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	productService := services.NewProductService(userRepoMock, productRepoMock, orderRepoMock)

	t.Run("success", func(t *testing.T) {
		expectedProducts := []entities.Product{
			{
				ID:        "1",
				Title:     "Product A",
				Price:     100,
				Detail:    "Description A",
				CreatedBy: "123",
			},
			{
				ID:        "2",
				Title:     "Product B",
				Price:     200,
				Detail:    "Description B",
				CreatedBy: "123",
			},
		}
		productRepoMock.On("Find", mock.Anything).Return(expectedProducts, nil)
		products, _ := productService.GetAll(ctx)
		assert.Equal(t, expectedProducts, products)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_product_find", func(t *testing.T) {
		productRepoMock.On("Find", mock.Anything).Return(nil, errs.ErrProductNotFound)
		_, err := productService.GetAll(ctx)
		assert.EqualError(t, err, errs.ErrProductNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})
}

func TestEditOneProduct(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	productService := services.NewProductService(userRepoMock, productRepoMock, orderRepoMock)

	t.Run("success", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		newProduct := &entities.Product{ID: "4", Title: "apple", Detail: "just a normal sweet apple"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, expectedProduct.ID).Return(expectedProduct, nil)
		productRepoMock.On("EnsureUpdateProduct", mock.Anything, expectedProduct, newProduct).Return(newProduct, nil)
		productRepoMock.On("UpdateOne", mock.Anything, newProduct, newProduct.ID).Return(nil)

		err := productService.EditOne(ctx, newProduct, newProduct.ID, "example@example.com")
		assert.Equal(t, err, nil)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_editOne_findByEmail", func(t *testing.T) {
		product := &entities.Product{ID: "123", Title: "test", Detail: "just test details"}
		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(nil, errs.ErrUserNotFound)
		err := productService.EditOne(ctx, product, product.ID, "example@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_editOne_findById", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		product := &entities.Product{ID: "4", Title: "test", Detail: "just test details"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, product.ID).Return(nil, errs.ErrProductNotFound)

		err := productService.EditOne(ctx, product, product.ID, "example@example.com")
		assert.EqualError(t, err, errs.ErrProductNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil

	})

	t.Run("error_editOne_forbidden", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "123"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, expectedProduct.ID).Return(expectedProduct, nil)

		err := productService.EditOne(ctx, expectedProduct, expectedProduct.ID, "example@example.com")
		assert.EqualError(t, err, errs.ErrForbidden.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_editOne_ensure", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		newProduct := &entities.Product{ID: "4", Title: "m", Detail: "mm"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, expectedProduct.ID).Return(expectedProduct, nil)
		productRepoMock.On("EnsureUpdateProduct", mock.Anything, expectedProduct, newProduct).Return(newProduct, nil)

		err := productService.EditOne(ctx, newProduct, newProduct.ID, "example@example.com")
		assert.EqualError(t, err, errs.ErrInValidProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_editOne_updateOne", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		newProduct := &entities.Product{ID: "4", Title: "apple", Detail: "just a normal sweet apple"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, expectedProduct.ID).Return(expectedProduct, nil)
		productRepoMock.On("EnsureUpdateProduct", mock.Anything, expectedProduct, newProduct).Return(newProduct, nil)
		productRepoMock.On("UpdateOne", mock.Anything, newProduct, newProduct.ID).Return(errs.ErrUpdateProduct)

		err := productService.EditOne(ctx, newProduct, newProduct.ID, "example@example.com")
		assert.EqualError(t, err, errs.ErrUpdateProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

}

func TestDropOneProduct(t *testing.T) {
	ctx := context.Background()

	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()

	productService := services.NewProductService(userRepoMock, productRepoMock, orderRepoMock)

	t.Run("error_dropOne_findByEmail", func(t *testing.T) {
		userRepoMock.On("FindByEmail", ctx, "example@example.com").Return(nil, errs.ErrUserNotFound)
		err := productService.DropOne(ctx, "2", "example@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_dropOne_findById", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, "4").Return(nil, errs.ErrProductNotFound)

		err := productService.DropOne(ctx, "4", "example@example.com")

		assert.EqualError(t, err, errs.ErrProductNotFound.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_dropOne_createdByNotMatch", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "4"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, "4").Return(expectedProduct, nil)

		err := productService.DropOne(ctx, "4", "example@example.com")

		assert.EqualError(t, err, errs.ErrForbidden.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
	})

	t.Run("error_dropOne_deleteAllOrderByProduct", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", mock.Anything, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", mock.Anything, "4").Return(expectedProduct, nil)
		orderRepoMock.On("DeleteAllOrderByProduct", mock.Anything, expectedProduct.ID).Return(errs.ErrDeleteAllOrderByProduct)

		err := productService.DropOne(ctx, "4", "example@example.com")

		assert.EqualError(t, err, errs.ErrDeleteAllOrderByProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_dropOne_deleteOne", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", ctx, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", ctx, expectedProduct.ID).Return(expectedProduct, nil)
		orderRepoMock.On("DeleteAllOrderByProduct", ctx, expectedProduct.ID).Return(nil)
		productRepoMock.On("DeleteOne", ctx, expectedProduct.ID).Return(errs.ErrDeleteProduct)

		err := productService.DropOne(ctx, "4", "example@example.com")

		assert.EqualError(t, err, errs.ErrDeleteProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		expectedUser := &entities.User{ID: "12", Email: "example@example.com"}
		expectedProduct := &entities.Product{ID: "4", Title: "test", Detail: "just test details", CreatedBy: "12"}

		userRepoMock.On("FindByEmail", ctx, "example@example.com").Return(expectedUser, nil)
		productRepoMock.On("FindById", ctx, expectedProduct.ID).Return(expectedProduct, nil)
		orderRepoMock.On("DeleteAllOrderByProduct", ctx, expectedProduct.ID).Return(nil)
		productRepoMock.On("DeleteOne", ctx, expectedProduct.ID).Return(nil)

		err := productService.DropOne(ctx, "4", "example@example.com")

		assert.Nil(t, err)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})
}
