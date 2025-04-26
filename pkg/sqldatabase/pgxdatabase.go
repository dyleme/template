package sqldatabase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGX(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("new: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return conn, nil
}
