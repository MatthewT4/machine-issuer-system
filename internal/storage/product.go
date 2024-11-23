package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"machineIssuerSystem/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *PgStorage) GetAvailableServers(ctx context.Context) ([]model.Server, error) {
	sql := "SELECT * FROM servers WHERE rent_by is null"

	rows, err := p.connections.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("queryex: %w", err)
	}
	servers, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Server])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.ErrNotFound{}
		}

		return nil, err
	}

	return servers, nil
}

func (p *PgStorage) GetMyServers(ctx context.Context, userID uuid.UUID) ([]model.Server, error) {
	sql := "SELECT * FROM servers WHERE rent_by = $1"

	rows, err := p.connections.Query(ctx, sql, userID)
	if err != nil {
		return nil, fmt.Errorf("queryex: %w", err)
	}
	servers, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Server])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &model.ErrNotFound{}
		}

		return nil, err
	}

	return servers, nil
}

func (p *PgStorage) GetServer(ctx context.Context, serverID uuid.UUID) (model.Server, error) {
	sql := "SELECT * FROM servers WHERE id = $1"

	rows, err := p.connections.Query(ctx, sql, serverID)
	if err != nil {
		return model.Server{}, fmt.Errorf("queryex: %w", err)
	}
	server, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Server])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Server{}, &model.ErrNotFound{}
		}

		return model.Server{}, err
	}

	return server, nil
}

func (p *PgStorage) RentServer(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error {
	sql := "UPDATE servers SET rent_by = $1 WHERE id = $2"

	res, err := p.connections.Exec(ctx, sql, userID, serverID)
	if err != nil {
		return fmt.Errorf("RentServer exec: %w", err)
	}
	p.logger.Info("RentServer", slog.Int64("count_updated", res.RowsAffected()))

	return nil
}

func (p *PgStorage) UnRentServer(ctx context.Context, serverID uuid.UUID) error {
	sql := "UPDATE servers SET rent_by = Null WHERE id = $1"

	res, err := p.connections.Exec(ctx, sql, serverID)
	if err != nil {
		return fmt.Errorf("UnRentServer exec: %w", err)
	}
	p.logger.Info("UnRentServer", slog.Int64("count_updated", res.RowsAffected()))

	return nil
}

func (p *PgStorage) GetServerIp(ctx context.Context, serverID uuid.UUID) (string, error) {
	sql := "SELECT ip FROM servers WHERE id = $1"

	var ip string
	err := p.connections.QueryRow(ctx, sql, serverID).Scan(&ip)
	if err != nil {
		return "", fmt.Errorf("GetServerIp: %w", err)
	}

	return ip, nil
}
