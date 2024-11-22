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
	GetAvailableServers(ctx context.Context) ([]model.Server, error)
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

func (c *Core) GetAvailableServers(ctx context.Context) ([]model.Server, error) {
	servers, err := c.storage.GetAvailableServers(ctx)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return nil, err
		}
		return nil, fmt.Errorf("GetAvailableServers from db: %w", err)
	}
	return servers, err
}
