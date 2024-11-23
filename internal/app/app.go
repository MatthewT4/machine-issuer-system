package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"

	appConfig "machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/controller"
	appCore "machineIssuerSystem/internal/core"
	"machineIssuerSystem/internal/storage"
	"machineIssuerSystem/internal/workers"
)

type Application struct {
	logger   *slog.Logger
	storage  *storage.PgStorage
	core     *appCore.Core
	api      *controller.Controller
	registry *prometheus.Registry
	cfg      appConfig.Config
}

func NewApplication(logger *slog.Logger) (*Application, error) {
	config, err := appConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	logger.Info("config loaded successfully")

	registry := prometheus.NewRegistry()

	pgStorage := storage.NewPgStorage(config.DbURL, logger)

	core := appCore.NewCore(pgStorage, logger, config)

	api := controller.NewController(core, config.ApiServerPort, logger, config, registry)

	return &Application{
		logger:   logger,
		storage:  pgStorage,
		core:     core,
		api:      api,
		registry: registry,
		cfg:      config,
	}, nil
}

func (p *Application) Start() error {
	p.logger.Info("Starting app")

	err := p.storage.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("connection to database: %w", err)
	}

	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricWorker := workers.NewMetric(p.cfg, p.logger, p.core)
	go metricWorker.Start(stopCtx)

	return p.api.Start()
}
