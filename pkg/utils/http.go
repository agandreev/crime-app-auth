package utils

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/agandreev/crime-app-auth/internal/domain"
	"github.com/agandreev/crime-app-auth/internal/repository"
)

func StatusCode(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, repository.ErrDBError):
		return http.StatusInternalServerError
	case errors.Is(err, repository.ErrExistedUser) ||
		errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) ||
		errors.Is(err, domain.ErrIncorrectUser):
		return http.StatusBadRequest
	case errors.Is(err, repository.ErrNoUser):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
