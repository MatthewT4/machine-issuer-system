package controller

import (
	"errors"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"log/slog"
	"machineIssuerSystem/internal/core"
	"machineIssuerSystem/internal/model"
	"net/http"
)

type productHandlers struct {
	core *core.Core

	logger *slog.Logger
}

func newProductHandlers(core *core.Core, logger *slog.Logger) *productHandlers {
	return &productHandlers{
		core: core,

		logger: logger,
	}
}

// (POST /rent/{server_id})
func (p *productHandlers) RentServer(ctx echo.Context, serverId openapi_types.UUID) error {
	return nil
}

// (GET /servers/available)
func (p *productHandlers) GetAvailableServers(ctx echo.Context) error {
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

func (p *productHandlers) GetProduct(ctx echo.Context, productId openapi_types.UUID) error {
	p.logger.Info("Get product request", slog.Any("productId", productId))

	product, err := p.core.GetProduct(ctx.Request().Context(), productId)
	if err != nil {
		var errNotFound *model.ErrNotFound
		switch {
		case errors.As(err, &errNotFound):
			p.logger.Debug("Product not found", slog.Any("productId", productId))
			return echo.NewHTTPError(http.StatusNotFound, "Produвct not found")
		default:
			p.logger.Error("Get product unknown error", slog.Any("productId", productId), slog.Any("error", err))
			return echo.ErrInternalServerError
		}
	}

	return ctx.JSON(http.StatusOK, product)
}
