package workers

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
)

type Metric struct {
	cfg    config.Config
	logger *slog.Logger

	core *core.Core

	serverCount prometheus.Gauge
}

func NewMetric(cfg config.Config, logger *slog.Logger, core *core.Core) *Metric {
	serverCountGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "servers_count",
		Help: "Number of servers",
	})
	prometheus.MustRegister(serverCountGauge)

	return &Metric{
		cfg:         cfg,
		logger:      logger,
		core:        core,
		serverCount: serverCountGauge,
	}
}

func (m *Metric) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(2) * time.Second)

	for {
		select {
		case <-ticker.C:
			m.logger.Info("Metric tick")

			servers, err := m.core.GetAvailableServers(context.Background())
			if err != nil {
				m.logger.Error("Error getting available servers", err)
			}

			m.logger.Info("setting metric server count", len(servers))
			m.serverCount.Set(float64(len(servers)))
		case <-ctx.Done():
			m.logger.Info("Worker stopping")
		}
	}
}
