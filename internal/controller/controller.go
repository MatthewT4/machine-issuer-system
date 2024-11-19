package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log/slog"
	"machineIssuerSystem/api"
	"machineIssuerSystem/internal/core"
)

type Controller struct {
	core   *core.Core
	server *echo.Echo
	port   uint16
	logger *slog.Logger
}

func NewController(core *core.Core, port uint16, logger *slog.Logger) *Controller {
	e := echo.New()

	productHandlers := newProductHandlers(core, logger)
	api.RegisterHandlers(e, productHandlers)

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
