package controller

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
	"machineIssuerSystem/internal/model"
)

var userID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type handlers struct {
	core *core.Core

	logger *slog.Logger
	cfg    config.Config
}

func newHandlers(core *core.Core, logger *slog.Logger, cfg config.Config) *handlers {
	return &handlers{
		core: core,

		logger: logger,
		cfg:    cfg,
	}
}

// (POST /rent/{server_id})
func (p *handlers) RentServer(ctx echo.Context, serverId openapi_types.UUID) error {
	err := p.core.RentServer(
		ctx.Request().Context(),
		userID,
		serverId,
	)
	return p.convertCoreErrorToResponse(err)
}

func (p *handlers) UnRentServer(ctx echo.Context, serverId openapi_types.UUID) error {
	p.logger.Debug("handle UnRentServer", slog.Any("server_id", serverId))
	err := p.core.UnRentServer(
		ctx.Request().Context(),
		serverId,
	)
	return p.convertCoreErrorToResponse(err)
}

// (GET /servers/available)
func (p *handlers) GetAvailableServers(ctx echo.Context) error {
	servers, err := p.core.GetAvailableServers(ctx.Request().Context())
	if err != nil {
		var errNotFound *model.ErrNotFound
		switch {
		case errors.As(err, &errNotFound):
			p.logger.Debug("servers not found")
			return echo.NewHTTPError(http.StatusNotFound, "Servers not found")
		default:
			p.logger.Error("GetAvailableServers unknown error", slog.Any("error", err))
			return echo.ErrInternalServerError
		}
	}
	return ctx.JSON(http.StatusOK, servers)
}

func (p *handlers) GetMyServers(ctx echo.Context) error {
	reqUserID := userID
	servers, err := p.core.GetMyServers(ctx.Request().Context(), reqUserID)
	if err != nil {
		var errNotFound *model.ErrNotFound
		switch {
		case errors.As(err, &errNotFound):
			p.logger.Debug("servers not found")
			return echo.NewHTTPError(http.StatusNotFound, "Servers not found")
		default:
			p.logger.Error("GetAvailableServers unknown error", slog.Any("error", err))
			return echo.ErrInternalServerError
		}
	}
	return ctx.JSON(http.StatusOK, servers)
}

func (p *handlers) convertCoreErrorToResponse(err error) error {
	if err == nil {
		return nil
	}
	var errBadRequest *model.ErrBadRequest
	var errNotFound *model.ErrNotFound
	var errInternal *model.ErrInternal

	switch {
	case errors.As(err, &errBadRequest):
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	case errors.As(err, &errNotFound):
		return echo.NewHTTPError(http.StatusNotFound, "Not found")
	case errors.As(err, &errInternal):
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	default:
		p.logger.Error("Unknown error", slog.Any("error", err))
		return echo.ErrInternalServerError
	}
}
