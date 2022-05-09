package domain

import (
	"errors"
	"fmt"
)

var (
	ErrIncorrectUser = errors.New("user values are incorrect")
)

type User struct {
	Login    string  `json:"login"`
	Password *string `json:"password,omitempty"`
}

func (u User) Validate() error {
	if len(u.Login) == 0 || len(u.Login) > 16 {
		return fmt.Errorf(
			"login should contain more than 0 and less than 16 symbols, err: %w",
			ErrIncorrectUser,
		)
	}
	if len(*u.Password) == 0 || len(*u.Password) > 16 {
		return fmt.Errorf(
			"password should contain more than 0 and less than 16 symbols, err: %w",
			ErrIncorrectUser,
		)
	}
	return nil
}
