package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"machineIssuerSystem/internal/model"
)

type Storage interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (model.Product, error)
}

type Core struct {
	storage Storage

	logger *slog.Logger
}

func NewCore(storage Storage, logger *slog.Logger) *Core {
	return &Core{
		storage: storage,

		logger: logger,
	}
}

func (c *Core) GetProduct(ctx context.Context, productID uuid.UUID) (model.Product, error) {
	product, err := c.storage.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return model.Product{}, err
		}
		return model.Product{}, fmt.Errorf("get product: %w", err)
	}
	return product, err
}
