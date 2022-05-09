package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/agandreev/crime-app-auth/internal/domain"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	selectUser = "SELECT login, password FROM users WHERE login=$1"
	insertUser = "INSERT INTO users (login, password) VALUES ($1, $2) ON CONFLICT DO NOTHING"
)

var (
	ErrExistedUser = errors.New("user is already existed")
	ErrNoUser      = errors.New("such user doesn't exist")
	ErrDBError     = errors.New("db error is met")
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) User(ctx context.Context, login string) (*domain.User, error) {
	user := domain.User{}
	if err := pgxscan.Get(ctx, u.db, &user, selectUser, login); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("can't get user from db, err: %w", ErrNoUser)
		}
		return nil, fmt.Errorf("can't get user from db, err: %w", ErrDBError)
	}
	return &user, nil
}

func (u UserRepository) AddUser(ctx context.Context, user domain.User) error {
	if selectedUser, _ := u.User(ctx, user.Login); selectedUser != nil {
		return fmt.Errorf("can't add user to db, err: %w", ErrExistedUser)
	}
	if _, err := u.db.Exec(ctx, insertUser, user.Login, user.Password); err != nil {
		return fmt.Errorf("can't add user to db, err: %w", ErrDBError)
	}
	return nil
}

func (u UserRepository) Close() {
	u.db.Close()
}
