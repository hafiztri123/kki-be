package config

import (
	"context"
	"fmt"

	"github.com/hafiztri123/kki-be/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(ctx context.Context) (*pgxpool.Pool, error) {

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		utils.GetEnv("DB_USERNAME"),
		utils.GetEnv("DB_PASSWORD"),
		utils.GetEnv("DB_HOST"),
		utils.GetEnv("DB_PORT"),
		utils.GetEnv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil

}
