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
	uptime      *prometheus.GaugeVec
	cpuUsage    *prometheus.GaugeVec
	ramUsage    *prometheus.GaugeVec
	memUsage    *prometheus.GaugeVec
}

func NewMetric(cfg config.Config, logger *slog.Logger, core *core.Core) *Metric {
	serverCountGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "servers_count",
		Help: "Number of servers",
	})

	uptimeGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "vm_uptime",
		Help: "Uptime of the resource",
	},
		[]string{"vm_title"})

	cpuUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "CPU usage",
	},
		[]string{"vm_title"})

	ramUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ram_usage",
		Help: "RAM usage",
	},
		[]string{"vm_title"})

	memUsageGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mem_usage",
		Help: "MEM usage",
	},
		[]string{"vm_title"})

	prometheus.MustRegister(serverCountGauge)
	prometheus.MustRegister(uptimeGauge)
	prometheus.MustRegister(cpuUsageGauge)
	prometheus.MustRegister(ramUsageGauge)
	prometheus.MustRegister(memUsageGauge)

	return &Metric{
		cfg:         cfg,
		logger:      logger,
		core:        core,
		serverCount: serverCountGauge,
		uptime:      uptimeGauge,
		cpuUsage:    cpuUsageGauge,
		ramUsage:    ramUsageGauge,
		memUsage:    memUsageGauge,
	}
}

func (m *Metric) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)

	for {
		select {
		case <-ticker.C:
			m.logger.Info("Metric tick")

			servers, err := m.core.GetAvailableServers(context.Background())
			if err != nil {
				m.logger.Error("Error getting available servers", err)
			}
			m.serverCount.Set(float64(len(servers)))

			for _, server := range servers {
				metric, err := m.core.GetMetrics(ctx, server.ID)
				if err != nil {
					m.logger.Error("Error getting metrics", err)
					continue
				}
				m.logger.Info("metric", metric)

				m.uptime.WithLabelValues(server.Title).Set(float64(metric.Uptime))
				m.cpuUsage.WithLabelValues(server.Title).Set(metric.CPU)
				m.ramUsage.WithLabelValues(server.Title).Set(metric.RAM)
				m.memUsage.WithLabelValues(server.Title).Set(float64(metric.Memory))
			}
		case <-ctx.Done():
			m.logger.Info("Worker stopping")
		}
	}
}
