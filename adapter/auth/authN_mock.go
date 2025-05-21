package authAdapter

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type authenServiceMock struct {
	mock.Mock
}

func NewAuthenServiceMock() *authenServiceMock {
	return &authenServiceMock{}
}

func (a *authenServiceMock) CreateToken(email string) (string, error) {
	args := a.Called(email)
	if args.Get(0) == "" {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

func (a *authenServiceMock) VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	args := a.Called(tokenStr)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
