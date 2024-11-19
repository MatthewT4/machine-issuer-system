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

type ProductorApp struct {
	logger  *slog.Logger
	storage *storage.PgStorage
	core    *appCore.Core
	api     *controller.Controller
}

func NewProductorApp(logger *slog.Logger) (*ProductorApp, error) {
	config, err := appConfig.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("get config: %w", err)
	}

	pgStorage := storage.NewPgStorage(config.DbURL)

	core := appCore.NewCore(pgStorage, logger)
	api := controller.NewController(core, config.ApiServerPort, logger)

	return &ProductorApp{
		logger:  logger,
		storage: pgStorage,
		core:    core,
		api:     api,
	}, nil
}

func (p *ProductorApp) Start() error {
	p.logger.Info("Starting app")
	err := p.storage.Connect(context.Background())
	if err != nil {
		return fmt.Errorf("connection to database: %w", err)
	}

	return p.api.Start()
}
