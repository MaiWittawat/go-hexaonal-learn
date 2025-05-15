package errs

import "errors"

var (
	// General Error
	ErrForbidden        = errors.New("no permission")
	ErrEmailAlreadyUsed = errors.New("email already used")
	ErrInvalidInput     = errors.New("invalid input")
	ErrUnauthorize      = errors.New("unauthorize")
	ErrCreatToken       = errors.New("fail to create token")
	ErrVerifyToken      = errors.New("fail to verify token")
	ErrLogin            = errors.New("fail to login")
	ErrPassNotMatch     = errors.New("password or ")

	// User Error
	ErrUserNotFound = errors.New("user not found")
	ErrInValidUser  = errors.New("invalid user")
	ErrHashPassword = errors.New("fail to hash password")
	ErrSaveUser     = errors.New("fail to save user")
	ErrUpdateUser   = errors.New("fail to update user")
	ErrDeleteUser   = errors.New("fail to delete user")

	// Product Error
	ErrProductNotFound  = errors.New("product not found")
	ErrInValidProduct   = errors.New("invalid product")
	ErrSaveProduct      = errors.New("fail to save product")
	ErrUpdateProduct    = errors.New("fail to update product")
	ErrDeleteProduct    = errors.New("fail to delete product")
	ErrDeleteAllProduct = errors.New("fail to delete all product")

	// Order Error
	ErrOrderNotFound           = errors.New("order not found")
	ErrDeleteAllOrderByProduct = errors.New("fail to delete order by product_id")
	ErrDeleteAllOrderByUser    = errors.New("fail to delete order by user_id")
	ErrSaveOrder               = errors.New("fail to save order")
	ErrUpdateOrder             = errors.New("fail to update order")
	ErrDeleteOrder             = errors.New("fail to delete order")
)
