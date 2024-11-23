package controller

import (
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"machineIssuerSystem/api"
	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
)

type Controller struct {
	core     *core.Core
	server   *echo.Echo
	port     uint16
	logger   *slog.Logger
	registry *prometheus.Registry
}

func NewController(core *core.Core, port uint16, logger *slog.Logger, cfg config.Config, registry *prometheus.Registry) *Controller {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4200"}, // Замените на разрешённые домены
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodOptions},
		AllowHeaders:     []string{"Refer", echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	productHandlers := newHandlers(core, logger, cfg)
	api.RegisterHandlers(e, productHandlers)

	//e.Use(productHandlers.AuthMiddleware)
	//e.Use(productHandlers.PermissionMiddleware)
	return &Controller{
		core:     core,
		server:   e,
		port:     port,
		logger:   logger,
		registry: registry,
	}
}

func (c *Controller) Start() error {
	address := fmt.Sprintf("[::]:%v", c.port)
	c.logger.Info("Server starting", slog.String("address", address))

	c.server.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return c.server.Start(fmt.Sprintf("[::]:%d", c.port))
}
