package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/agandreev/crime-app-auth/internal/domain"
	"github.com/agandreev/crime-app-auth/internal/handlers"
	"github.com/agandreev/crime-app-auth/pkg/config"
	"github.com/golang-jwt/jwt"
)

type UserRepository interface {
	User(ctx context.Context, login string) (*domain.User, error)
	AddUser(ctx context.Context, user domain.User) error
}

type UserService struct {
	users UserRepository

	timeout   time.Duration
	jwtConfig config.JWTConfig
}

func NewUserService(rep UserRepository, jwtConfig config.JWTConfig) handlers.UserService {
	return &UserService{
		users:     rep,
		jwtConfig: jwtConfig,
	}
}

func (u UserService) SignIn(ctx context.Context, user domain.User) (string, error) {
	storedUser, err := u.users.User(ctx, user.Login)
	if err != nil {
		return "", fmt.Errorf("can't sign in 'cause of user existence, err: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*storedUser.Password), []byte(*user.Password)); err != nil {
		return "", fmt.Errorf("can't sign in 'cause of password match, err: %w", err)
	}

	token, err := generateJWT(
		user.Login,
		u.jwtConfig.Expiration*time.Minute,
		u.jwtConfig.Secret,
	)
	if err != nil {
		return "", fmt.Errorf("can't sign in 'cause of JWT creation, err: %w", err)
	}
	return token, nil
}

func (u UserService) SignUp(ctx context.Context, user domain.User) error {
	hashedPassword, err := hash(*user.Password)
	if err != nil {
		return fmt.Errorf("can't sign up 'cause of hashing, err: %w", err)
	}
	user.Password = &hashedPassword

	if err = u.users.AddUser(ctx, user); err != nil {
		return fmt.Errorf("can't sign up 'cause of, err: %w", err)
	}
	return nil
}

func hash(password string) (string, error) {
	hashBinary, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hashBinary), nil
}

type jwtCrimeClaims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func generateJWT(login string, timeout time.Duration, secret string) (string, error) {
	claims := jwtCrimeClaims{
		Login: login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeout).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}
