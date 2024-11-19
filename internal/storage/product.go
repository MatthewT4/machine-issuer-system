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
