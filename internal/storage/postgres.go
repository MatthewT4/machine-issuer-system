package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type PgStorage struct {
	connections *pgxpool.Pool
	dbURL       string

	logger *slog.Logger
}

func NewPgStorage(dbURL string) *PgStorage {
	return &PgStorage{dbURL: dbURL}
}

func (p *PgStorage) Connect(ctx context.Context) error {

	config, err := pgxpool.ParseConfig(p.dbURL)
	if err != nil {
		return fmt.Errorf("parce config: %w", err)
	}

	connect, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("create connect: %v", err)
	}
	err = connect.Ping(ctx)
	if err != nil {
		return fmt.Errorf("ping connect: %w", err)
	}

	p.connections = connect
	return nil
}
