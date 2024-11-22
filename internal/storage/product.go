package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"machineIssuerSystem/internal/model"
)

func (p *PgStorage) GetProduct(ctx context.Context, productID uuid.UUID) (model.Product, error) {
	sql := "SELECT * FROM product WHERE id = $1"

	rows, err := p.connections.Query(ctx, sql, productID)
	if err != nil {
		return model.Product{}, fmt.Errorf("queryex: %w", err)
	}
	product, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Product])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Product{}, &model.ErrNotFound{}
		}

		return model.Product{}, err
	}

	return convertProductFromDB(product), nil
}

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

func (p *PgStorage) RentServer(ctx context.Context, serverId uuid.UUID, userId uuid.UUID) ([]model.Server, error) {
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
