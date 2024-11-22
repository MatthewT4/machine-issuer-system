package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"machineIssuerSystem/internal/model"
)

type Storage interface {
	GetProduct(ctx context.Context, productID uuid.UUID) (model.Product, error)
	GetAvailableServers(ctx context.Context) ([]model.Server, error)
	GetServer(ctx context.Context, serverID uuid.UUID) (model.Server, error)
	RentServer(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error
	UnRentServer(ctx context.Context, serverID uuid.UUID) error

	CreateUser(ctx context.Context, user model.User) (result model.User, err error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
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

func (c *Core) RentServer(ctx context.Context, userID uuid.UUID, serverID uuid.UUID) error {
	c.logger.Info("fdfd", slog.Any("dfd", serverID))
	server, err := c.storage.GetServer(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return &model.ErrNotFound{}
		}
		return &model.ErrBadRequest{}
	}
	if server.RentBy != nil {
		return &model.ErrBadRequest{}
	}
	err = c.storage.RentServer(ctx, serverID, userID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return err
		}
		c.logger.Error("rent fail", slog.Any("server", server), slog.Any("user", userID), slog.Any("error", err))
	}
	c.logger.Debug("rent ok", slog.Any("server", server), slog.Any("user", userID))
	return nil
}

func (c *Core) UnRentServer(ctx context.Context, serverID uuid.UUID) error {
	server, err := c.storage.GetServer(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return &model.ErrNotFound{}
		}
		return &model.ErrBadRequest{}
	}
	if server.RentBy == nil {
		return &model.ErrBadRequest{}
	}
	err = c.storage.UnRentServer(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return err
		}
		c.logger.Error("unrent fail", slog.Any("server", server), slog.Any("error", err))
	}
	c.logger.Debug("unrent ok", slog.Any("server", server))
	return nil
}
