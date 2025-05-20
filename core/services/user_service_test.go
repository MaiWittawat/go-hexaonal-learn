package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authAdapter "github.com/wittawat/go-hex/adapter/auth"
	orderAdapter "github.com/wittawat/go-hex/adapter/order/outbound"
	productAdapter "github.com/wittawat/go-hex/adapter/product/outbound"
	userAdapter "github.com/wittawat/go-hex/adapter/user/outbound"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/services"
	"github.com/wittawat/go-hex/utils/errs"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("error_isValidUser", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "u", Username: "u", Password: "u"}
		userRepoMock.On("IsValidUser", user).Return(nil, errs.ErrInValidUser)

		err := userService.Save(ctx, user, "user")
		assert.EqualError(t, err, errs.ErrInValidUser.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_save", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "user@example.com", Username: "user", Password: "user"}
		userRepoMock.On("IsValidUser", user).Return(nil)
		userRepoMock.On("Save", mock.Anything, user).Return(errs.ErrSaveUser)

		err := userService.Save(ctx, user, "user")
		assert.EqualError(t, err, errs.ErrSaveUser.Error())

		userRepoMock.ExpectedCalls = nil
	})
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("error_user_findById", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com", Username: "mai", Password: "mai"}

		userRepoMock.On("FindById", ctx, user.ID).Return(nil, errs.ErrUserNotFound)
		_, _, err := userService.UpdateOne(ctx, user, user.ID, "testToken")

		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_validUser", func(t *testing.T) {
		user := &entities.User{ID: "1", Email: "mai@example.com", Username: "mai", Password: "mai"}

		userRepoMock.On("FindById", ctx, user.ID).Return(user, nil)
		_, _, err := userService.UpdateOne(ctx, user, user.ID, "testToken")

		assert.EqualError(t, err, errs.ErrInValidUser.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_updateOne", func(t *testing.T) {
		tokenRepoMock := authAdapter.NewAuthNServiceImpMock()
		mockToken := "mocked.token.value"
		user := &entities.User{
			ID:       "1",
			Email:    "mai@example.com",
			Username: "mai",
			Password: "maimai",
		}

		userRepoMock.On("FindById", ctx, user.ID).Return(user, nil)
		tokenRepoMock.On("CreateToken", user.Email).Return(mockToken, nil)

		userRepoMock.On("UpdateOne", ctx, mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "1" && u.Email == "mai@example.com" && u.Username == "mai"
		}), "1").Return(errs.ErrUpdateUser)

		_, _, err := userService.UpdateOne(ctx, user, user.ID, "mai@example.com")

		assert.EqualError(t, err, errs.ErrUpdateUser.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		tokenRepoMock := authAdapter.NewAuthNServiceImpMock()
		mockToken := "mocked.token.value"
		user := &entities.User{
			ID:       "1",
			Email:    "mai@example.com",
			Username: "mai",
			Password: "maimai",
		}

		userRepoMock.On("FindById", ctx, user.ID).Return(user, nil)
		tokenRepoMock.On("CreateToken", user.Email).Return(mockToken, nil)

		userRepoMock.On("UpdateOne", ctx, mock.MatchedBy(func(u *entities.User) bool {
			return u.ID == "1" && u.Email == "mai@example.com" && u.Username == "mai"
		}), "1").Return(nil)

		_, _, err := userService.UpdateOne(ctx, user, user.ID, "mai@example.com")

		assert.Nil(t, err)

		userRepoMock.ExpectedCalls = nil
	})
}

func TestDeleteUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("error_delete_findById", func(t *testing.T) {
		expectedUser := &entities.User{
			ID:       "1",
			Email:    "mai@example.com",
			Username: "mai",
			Password: "maimai",
		}
		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(nil, errs.ErrUserNotFound)

		err := userService.DeleteOne(ctx, expectedUser.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_delete_emailNotMatch", func(t *testing.T) {
		expectedUser := &entities.User{
			ID:       "1",
			Email:    "mai@example.com",
			Username: "mai",
			Password: "maimai",
		}
		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(expectedUser, nil)

		err := userService.DeleteOne(ctx, expectedUser.ID, "test@example.com")
		assert.EqualError(t, err, errs.ErrForbidden.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_delete_deleteAllOrderByUser", func(t *testing.T) {
		expectedUser := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "password"}

		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(expectedUser, nil)
		orderRepoMock.On("DeleteAllOrderByUser", ctx, expectedUser.ID).Return(errs.ErrDeleteAllOrderByUser)

		err := userService.DeleteOne(ctx, expectedUser.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrDeleteAllOrderByUser.Error())

		userRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_delete_deleteAll", func(t *testing.T) {
		expectedUser := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "password"}
		productId := "1"

		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(expectedUser, nil)
		orderRepoMock.On("DeleteAllOrderByUser", ctx, expectedUser.ID).Return(nil)
		productRepoMock.On("DeleteAll", ctx, productId).Return(errs.ErrDeleteProduct)

		err := userService.DeleteOne(ctx, expectedUser.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrDeleteProduct.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("error_delete_deleteOne", func(t *testing.T) {
		expectedUser := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "password"}
		productId := "1"

		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(expectedUser, nil)
		orderRepoMock.On("DeleteAllOrderByUser", ctx, expectedUser.ID).Return(nil)
		productRepoMock.On("DeleteAll", ctx, productId).Return(nil)
		userRepoMock.On("DeleteOne", ctx, expectedUser.ID).Return(errs.ErrDeleteUser)

		err := userService.DeleteOne(ctx, expectedUser.ID, "mai@example.com")
		assert.EqualError(t, err, errs.ErrDeleteUser.Error())

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})

	t.Run("success", func(t *testing.T) {
		expectedUser := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "password"}
		productId := "1"

		userRepoMock.On("FindById", ctx, expectedUser.ID).Return(expectedUser, nil)
		orderRepoMock.On("DeleteAllOrderByUser", ctx, expectedUser.ID).Return(nil)
		productRepoMock.On("DeleteAll", ctx, productId).Return(nil)
		userRepoMock.On("DeleteOne", ctx, expectedUser.ID).Return(nil)

		err := userService.DeleteOne(ctx, expectedUser.ID, "mai@example.com")
		assert.Nil(t, err)

		userRepoMock.ExpectedCalls = nil
		productRepoMock.ExpectedCalls = nil
		orderRepoMock.ExpectedCalls = nil
	})
}

func TestFindUser(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("success", func(t *testing.T) {
		expectedValue := []entities.User{
			{ID: "1", Username: "user1", Email: "user1@example.com", Password: "user1"},
			{ID: "2", Username: "user2", Email: "user2@example.com", Password: "user2"},
			{ID: "3", Username: "user3", Email: "user3@example.com", Password: "user3"},
		}
		userRepoMock.On("Find", ctx).Return(expectedValue, nil)
		users, _ := userService.Find(ctx)
		assert.Equal(t, expectedValue, users)

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("success_user_findNil", func(t *testing.T) {
		expectedValue := []entities.User{}
		userRepoMock.On("Find", ctx).Return(expectedValue, nil)
		users, _ := userService.Find(ctx)
		assert.Equal(t, expectedValue, users)

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_find", func(t *testing.T) {
		userRepoMock.On("Find", ctx).Return(nil, errs.ErrUserNotFound)
		_, err := userService.Find(ctx)
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})
}

func TestFindUserById(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("success", func(t *testing.T) {
		expectedValue := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "mmai"}
		userRepoMock.On("FindById", ctx, expectedValue.ID).Return(expectedValue, nil)

		user, _ := userService.FindById(ctx, expectedValue.ID)
		assert.Equal(t, expectedValue, user)

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_findById", func(t *testing.T) {
		userRepoMock.On("FindById", ctx, "1").Return(nil, errs.ErrUserNotFound)
		_, err := userService.FindById(ctx, "1")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})
}

func TestFindUserByEmail(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("success", func(t *testing.T) {
		expectedValue := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "mmai"}
		userRepoMock.On("FindByEmail", ctx, expectedValue.Email).Return(expectedValue, nil)

		user, _ := userService.FindByEmail(ctx, expectedValue.Email)
		assert.Equal(t, expectedValue, user)

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_user_findByEmail", func(t *testing.T) {
		userRepoMock.On("FindByEmail", ctx, "mai@example.com").Return(nil, errs.ErrUserNotFound)
		_, err := userService.FindByEmail(ctx, "mai@example.com")
		assert.EqualError(t, err, errs.ErrUserNotFound.Error())

		userRepoMock.ExpectedCalls = nil
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	userRepoMock := userAdapter.NewUserRepositoryMock()
	productRepoMock := productAdapter.NewProductRepositoryMock()
	orderRepoMock := orderAdapter.NewOrderRepositoryMock()
	token := authAdapter.NewAuthNServiceImpl()

	userService := services.NewUserService(userRepoMock, productRepoMock, orderRepoMock, token)

	t.Run("error_login_findByeEmail", func(t *testing.T) {
		expectedValue := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "mmai"}
		userRepoMock.On("FindByEmail", ctx, expectedValue.Email).Return(nil, errs.ErrLogin)
		_, err := userService.Login(ctx, expectedValue)
		assert.EqualError(t, err, errs.ErrLogin.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_login_passwordNotMatch", func(t *testing.T) {
		data := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: "mai"}
		expectedValue := &entities.User{ID: "2", Username: "test", Email: "test@example.com", Password: "test"}
		userRepoMock.On("FindByEmail", ctx, data.Email).Return(expectedValue, nil)
		_, err := userService.Login(ctx, data)
		assert.EqualError(t, err, errs.ErrPassNotMatch.Error())

		userRepoMock.ExpectedCalls = nil
	})

	t.Run("error_login_invalid", func(t *testing.T) {
		password := "password"
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), 4)
		data := &entities.User{ID: "1", Username: "mai", Email: "mai@example.com", Password: password}

		expectedValue := &entities.User{ID: "2", Username: "test", Email: "mai@example.com", Password: string(hash)}

		userRepoMock.On("FindByEmail", ctx, data.Email).Return(expectedValue, nil)
		_, err := userService.Login(ctx, data)
		assert.EqualError(t, err, errs.ErrInvalidInput.Error())

		userRepoMock.ExpectedCalls = nil
	})

	// t.Run("success", func(t *testing.T) {
	// 	tokenRepoMock := authAdapter.NewAuthNServiceImpMock()
	// 	mockToken := "mocked.token.value"
	// 	user := &entities.User{
	// 		ID:       "1",
	// 		Email:    "mai@example.com",
	// 		Username: "mai",
	// 		Password: "maimai",
	// 	}

	// 	userRepoMock.On("FindById", ctx, user.ID).Return(user, nil)
	// 	tokenRepoMock.On("CreateToken", user.Email).Return(mockToken, nil)

	// 	token, _ := userService.Login(ctx, user)

	// 	assert.NotEmpty(t, token)
	// 	assert.Equal(t, mockToken, token)

	// 	userRepoMock.ExpectedCalls = nil
	// })

}
