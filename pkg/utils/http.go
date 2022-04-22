package utils

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/agandreev/crime-app-auth/internal/repository"
)

func StatusCode(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, repository.ErrDBError):
		return http.StatusInternalServerError
	case errors.Is(err, repository.ErrNoUser) ||
		errors.Is(err, repository.ErrExistedUser) ||
		errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
