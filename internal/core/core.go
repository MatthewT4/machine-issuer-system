package core

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/model"
	vm "machineIssuerSystem/internal/virtualmachine"
)

type Storage interface {
	GetAvailableServers(ctx context.Context) ([]model.Server, error)
	GetMyServers(ctx context.Context, userID uuid.UUID) ([]model.Server, error)
	GetServer(ctx context.Context, serverID uuid.UUID) (model.Server, error)
	RentServer(ctx context.Context, serverID uuid.UUID, userID uuid.UUID, rentUntil time.Time) error
	UnRentServer(ctx context.Context, serverID uuid.UUID) error
	ExpiredServers(ctx context.Context) ([]model.Server, error)

	CreateUser(ctx context.Context, user model.User) (result model.User, err error)
	GetUserByUsername(ctx context.Context, username string) (result model.User, err error)

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

func (c *Core) RentServer(ctx context.Context, userID uuid.UUID, serverID uuid.UUID, bookingDays int64) error {
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

	rentUntil := time.Now().Add(time.Duration(bookingDays) * 24 * time.Hour)

	err = c.storage.RentServer(ctx, serverID, userID, rentUntil)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return err
		}
		c.logger.Error("rent fail", slog.Any("server", server), slog.Any("user", userID), slog.Any("error", err))

		return err
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
		return err
	}
	c.logger.Debug("unrent ok", slog.Any("server", server))
	return nil
}

func (c *Core) GetMetrics(ctx context.Context, serverID uuid.UUID) (response model.Metric, err error) {
	log := c.logger.With(
		slog.String("op", "core.GetMetrics"),
		slog.String("server_id", serverID.String()))
	ip, err := c.storage.GetServerIp(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return response, &model.ErrNotFound{}
		}
		return response, err
	}
	log.Info("ip", ip)
	session, err := vm.CreateConnection(ip, c.cfg.SSHFilePath)
	if err != nil {
		log.Error("failed to create connection", err)
		return response, err
	}
	defer session.Close()
	metrics, err := vm.GetMetrics(
		session,
		[]string{vm.Uptime, vm.CPU, vm.RAM, vm.MEM},
	)
	if err != nil {
		log.Error("failed to request metrics", err)
		return response, err
	}
	return model.FromPkgToDomain(metrics), nil
}

func (c *Core) GetServer(ctx context.Context, serverID uuid.UUID) (model.Server, error) {
	server, err := c.storage.GetServer(ctx, serverID)
	return server, err
}

func (c *Core) RebootServer(ctx context.Context, serverID uuid.UUID) error {
	log := c.logger.With(
		slog.String("op", "core.RebootServer"),
		slog.String("server_id", serverID.String()))
	ip, err := c.storage.GetServerIp(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return &model.ErrNotFound{}
		}
		return err
	}
	log.Info("ip", ip)
	session, err := vm.CreateConnection(ip, c.cfg.SSHFilePath)
	if err != nil {
		log.Error("failed to create connection", err)
		return err
	}
	defer session.Close()
	commands := []string{vm.Reboot}
	err = vm.ExecuteCommandsOnVirtualMachine(session, commands)
	if err != nil {
		log.Error("failed to execute commands on virtual machine", err)
		return err
	}
	return nil
}

func (c *Core) CreateUserOnVm(ctx context.Context, serverID uuid.UUID) error {
	log := c.logger.With(
		slog.String("op", "core.GetMetrics"),
		slog.String("server_id", serverID.String()))
	ip, err := c.storage.GetServerIp(ctx, serverID)
	if err != nil {
		if errors.Is(err, &model.ErrNotFound{}) {
			return &model.ErrNotFound{}
		}
		return err
	}
	log.Info("ip", ip)
	session, err := vm.CreateConnection(ip, c.cfg.SSHFilePath)
	if err != nil {
		log.Error("failed to create connection", err)
		return err
	}
	defer session.Close()
	login := "user"
	password, err := generatePassword(10)
	if err != nil {
		log.Error("failed to generate password for vm user", err)
		return err
	}
	createUserCommands := []string{
		fmt.Sprintf(vm.CreateUser, login),
		fmt.Sprintf(vm.CreatePassword, password),
		fmt.Sprintf(vm.GiveRoot, login),
	}
	limitRootCommands := strings.Split(strings.TrimSpace(strings.ReplaceAll(vm.LimitRoot, "%s", "user")), "\n")
	commands := append(createUserCommands, limitRootCommands...)
	err = vm.ExecuteCommandsOnVirtualMachine(session, commands)
	if err != nil {
		log.Error("failed to execute commands on virtual machine", err)
		return err
	}
	return nil
}

func (c *Core) FetchExpiredServers(ctx context.Context) ([]model.Server, error) {
	log := c.logger.With(
		slog.String("op", "core.GetMetrics"))

	servers, err := c.storage.ExpiredServers(ctx)
	if err != nil {
		log.Error("failed to fetch expired servers", err)

		if errors.Is(err, &model.ErrNotFound{}) {
			return nil, &model.ErrNotFound{}
		}
		return nil, &model.ErrInternal{}
	}

	return servers, nil
}

func generatePassword(passwordLength int) (string, error) {
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?"
	password := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(allowedChars))))
		if err != nil {
			return "", err
		}
		password[i] = allowedChars[index.Int64()]
	}
	return string(password), nil
}
