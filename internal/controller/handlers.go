package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

const (
	defaultBookingDays = 7
)

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
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	request := model.RentServerRequest{}
	if err = json.Unmarshal(body, &request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if request.BookingDays == 0 {
		request.BookingDays = defaultBookingDays
	}

	err = p.core.RentServer(
		ctx.Request().Context(),
		userID,
		serverId,
		request.BookingDays,
	)
	fmt.Println(err)
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

func (p *handlers) GetServerMetrics(ctx echo.Context, serverId openapi_types.UUID) error {
	metrics, err := p.core.GetMetrics(ctx.Request().Context(), serverId)
	if err != nil {
		var errNotFound *model.ErrNotFound
		switch {
		case errors.As(err, &errNotFound):
			p.logger.Debug("server not found")
			return echo.NewHTTPError(http.StatusNotFound, "Server not found")
		default:
			p.logger.Error("GetServerMetrics unknown error", slog.Any("error", err))
			return echo.ErrInternalServerError
		}
	}

	return ctx.JSON(http.StatusOK, metrics)
}

func (h *handlers) GetServer(ctx echo.Context, serverId openapi_types.UUID) error {
	h.logger.Debug("handle GetServer", slog.Any("server_id", serverId))

	server, err := h.core.GetServer(ctx.Request().Context(), serverId)
	if err != nil {
		return h.convertCoreErrorToResponse(err)
	}
	return ctx.JSON(http.StatusOK, server)
}

func (p *handlers) RebootServer(ctx echo.Context, serverId openapi_types.UUID) error {
	err := p.core.RebootServer(
		ctx.Request().Context(),
		serverId,
	)
	return p.convertCoreErrorToResponse(err)
}

func (p *handlers) CreateUserOnVm(ctx echo.Context, serverId openapi_types.UUID) error {
	err := p.core.CreateUserOnVm(
		ctx.Request().Context(),
		serverId,
	)
	return p.convertCoreErrorToResponse(err)
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
