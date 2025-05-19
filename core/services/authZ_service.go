package services

import (
	"context"

	userPort "github.com/wittawat/go-hex/core/port/user"
)

type AuthorizeService struct {
	userRepo userPort.UserRepository
}

func NewAuthZServiceImpl(userRepo userPort.UserRepository) *AuthorizeService {
	return &AuthorizeService{userRepo: userRepo}
}

func (a *AuthorizeService) Authorize(email string, requiredRoles ...string) (bool, error) {
	user, err := a.userRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return false, err
	}

	for _, role := range requiredRoles {
		if role == user.Role {
			return true, nil
		}
	}
	return false, err
}
