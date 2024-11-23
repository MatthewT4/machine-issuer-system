package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/model"
	vm "machineIssuerSystem/internal/virtualmachine"
)

type Storage interface {
	GetAvailableServers(ctx context.Context) ([]model.Server, error)
	GetMyServers(ctx context.Context, userID uuid.UUID) ([]model.Server, error)
	GetServer(ctx context.Context, serverID uuid.UUID) (model.Server, error)
	RentServer(ctx context.Context, serverID uuid.UUID, userID uuid.UUID) error
	UnRentServer(ctx context.Context, serverID uuid.UUID) error

	CreateUser(ctx context.Context, user model.User) (result model.User, err error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	GetPermissionHandler(ctx context.Context, params model.GetPermissionHandlerRequest) (model.PermissionHandler, error)

	GetServerIp(ctx context.Context, serverID uuid.UUID) (string, error)
}

type Core struct {
	storage Storage

	logger *slog.Logger
	cfg    config.Config
}

func NewCore(storage Storage, logger *slog.Logger, cfg config.Config) *Core {
	return &Core{
		storage: storage,

		logger: logger,
		cfg:    cfg,
	}
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

func (c *Core) GetMyServers(ctx context.Context, userID uuid.UUID) ([]model.Server, error) {
	servers, err := c.storage.GetMyServers(ctx, userID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return nil, err
		}
		return nil, fmt.Errorf("GetMyServers from db: %w", err)
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

func (c *Core) GetMetrics(ctx context.Context, serverID uuid.UUID) error {
	ip, err := c.storage.GetServerIp(ctx, serverID)
	if err != nil {
		return &model.ErrBadRequest{}
	}
	session, err := vm.CreateConnection(ip)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	metrics, err := vm.RequestAndProcessMetrics(session)
	if err != nil {
		panic(err)
	}
	fmt.Println(metrics.Uptime, metrics.CPU, metrics.RAM, metrics.MEM)

	return nil
}
