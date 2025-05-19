package authAdapter

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type authNServiceImplMock struct {
	mock.Mock
}

func NewAuthNServiceImpMock() *authNServiceImplMock {
	return &authNServiceImplMock{}
}

func (a *authNServiceImplMock) CreateToken(email string) (string, error) {
	args := a.Called(email)
	if args.Get(0) == "" {
		return "", args.Error(1)
	}
	return args.Get(0).(string), nil
}

func (a *authNServiceImplMock) VerifyToken(tokenStr string) (jwt.MapClaims, error) {
	args := a.Called(tokenStr)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}
