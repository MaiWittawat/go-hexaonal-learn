package validator

import (
	"errors"
	"regexp"
	"time"

	"github.com/wittawat/go-hex/core/entities"
	"golang.org/x/crypto/bcrypt"
)

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func updateUser(oldUser *entities.User, newUser *entities.User) (*entities.User, error) {
	if newUser.Username == "" {
		newUser.Username = oldUser.Username
	}
	if newUser.Password == "" {
		newUser.Password = oldUser.Password
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
		if err != nil {
			return nil, err
		}
		newUser.Password = string(hash)
	}
	newUser.CreatedAt = oldUser.CreatedAt
	newUser.UpdatedAt = time.Now()
	return newUser, nil
}

func IsValidUser(user *entities.User) error {
	if len(user.Password) < 4 {
		return errors.New("invalid password more than 4 character")
	}
	if !isValidEmail(user.Email) {
		return errors.New("email is invalid")
	}
	return nil
}

func IsValidUpdateUser(oldUser *entities.User, newUser *entities.User, tokenEmail string) (*entities.User, error) {
	if !isValidEmail(newUser.Email) {
		return nil, errors.New("email is invalid")
	}
	if oldUser.Email != tokenEmail {
		return nil, errors.New("fail to update user wrong email")
	}
	user, err := updateUser(oldUser, newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}
