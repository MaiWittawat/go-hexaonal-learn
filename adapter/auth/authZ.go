// adapter/auth/authZ.go
package adapter

import (
	"errors"

	port "github.com/wittawat/go-hex/core/port/user"
)

type AuthZServiceImpl struct {
	userRepo port.UserOutbound
}

func NewAuthZServiceImpl(userRepo port.UserOutbound) *AuthZServiceImpl {
	return &AuthZServiceImpl{userRepo: userRepo}
}

func (a *AuthZServiceImpl) Authorize(email string, requiredRoles []string) (bool, error) {
	user, err := a.userRepo.FindByEmail(email)
	if err != nil {
		return false, err
	}

	for _, role := range requiredRoles {
		if role == user.Role {
			return true, nil
		}
	}
	return false, errors.New("forbidden")
}
