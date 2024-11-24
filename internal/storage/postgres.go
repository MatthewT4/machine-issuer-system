package storage

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgStorage struct {
	connections *pgxpool.Pool
	dbURL       string

	logger *slog.Logger
}

func NewPgStorage(dbURL string, logger *slog.Logger) *PgStorage {
	return &PgStorage{dbURL: dbURL, logger: logger}
}

func (p *PgStorage) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(p.dbURL)
	if err != nil {
		return nil, fmt.Errorf("parce config: %w", err)
	}

	connect, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create connect: %v", err)
	}

	err = connect.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping connect: %w", err)
	}

	p.connections = connect
	return connect, nil
}
