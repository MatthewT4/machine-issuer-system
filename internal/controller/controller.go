package controller

import (
	"fmt"

	"log/slog"

	"github.com/labstack/echo/v4"

	"machineIssuerSystem/api"
	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
)

type Controller struct {
	core   *core.Core
	server *echo.Echo
	port   uint16
	logger *slog.Logger
}

func NewController(core *core.Core, port uint16, logger *slog.Logger, cfg config.Config) *Controller {
	e := echo.New()

	productHandlers := newHandlers(core, logger, cfg)
	api.RegisterHandlers(e, productHandlers)

	e.Use(productHandlers.AuthMiddleware)
	e.Use(productHandlers.PermissionMiddleware)

	return &Controller{
		core:   core,
		server: e,
		port:   port,
		logger: logger,
	}
}

func (c *Controller) Start() error {
	address := fmt.Sprintf("[::]:%v", c.port)
	c.logger.Info("Server starting", slog.String("address", address))

	return c.server.Start(fmt.Sprintf("[::]:%d", c.port))
}
