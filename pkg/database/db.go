package database

import (
	"context"
	"fmt"

	"github.com/agandreev/crime-app-auth/pkg/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func Connect(ctx context.Context, config config.DatabaseConfig) (*pgxpool.Pool, error) {
	connection := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Port,
		config.Name,
	)
	return pgxpool.Connect(ctx, connection)
}
