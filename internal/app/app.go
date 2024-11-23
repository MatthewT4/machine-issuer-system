package app

import (
	"context"
	"fmt"
	"log/slog"

	appConfig "machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/controller"
	appCore "machineIssuerSystem/internal/core"
	"machineIssuerSystem/internal/storage"
)

type Application struct {
	logger  *slog.Logger
	storage *storage.PgStorage
	core    *appCore.Core
	api     *controller.Controller
}

func NewApplication(logger *slog.Logger) (*Application, error) {
	config, err := appConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	logger.Info("config loaded successfully")

	pgStorage := storage.NewPgStorage(config.DbURL, logger)

	core := appCore.NewCore(pgStorage, logger, config)

	api := controller.NewController(core, config.ApiServerPort, logger, config)

	return &Application{
		logger:  logger,
		storage: pgStorage,
		core:    core,
		api:     api,
	}, nil
}

func (p *Application) Start() error {
	p.logger.Info("Starting app")

	err := p.storage.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("connection to database: %w", err)
	}

	return p.api.Start()
}
