package workers

import (
	"log/slog"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
)

type Metric struct {
	cfg    *config.Config
	logger *slog.Logger

	core *core.Core
}

func NewMetric(cfg *config.Config, logger *slog.Logger, core *core.Core) *Metric {
	return &Metric{
		cfg:    cfg,
		logger: logger,
		core:   core,
	}
}
