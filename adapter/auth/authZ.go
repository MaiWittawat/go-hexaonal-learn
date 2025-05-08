// adapter/auth/authZ.go
package authAdapter

import (
	"errors"

	userPort "github.com/wittawat/go-hex/core/port/user"
)

type AuthZServiceImpl struct {
	userRepo userPort.UserRepository
}

func NewAuthZServiceImpl(userRepo userPort.UserRepository) *AuthZServiceImpl {
	return &AuthZServiceImpl{userRepo: userRepo}
}

func (a *AuthZServiceImpl) Authorize(email string, requiredRoles ...string) (bool, error) {
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
