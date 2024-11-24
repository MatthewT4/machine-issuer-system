package expired_server

import (
	"context"
	"log/slog"
	"time"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"

	"machineIssuerSystem/internal/config"
	"machineIssuerSystem/internal/core"
)

type Worker struct {
	cfg    config.Config
	logger *slog.Logger

	core      *core.Core
	txManager *manager.Manager
}

func NewWorker(cfg config.Config, logger *slog.Logger, core *core.Core, txManager *manager.Manager) *Worker {
	return &Worker{
		cfg:    cfg,
		logger: logger,

		core:      core,
		txManager: txManager,
	}
}

func (m *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)

	for {
		select {
		case <-ticker.C:
			m.logger.Info("Expired worker tick")

			servers, err := m.core.FetchExpiredServers(ctx)
			if err != nil {
				m.logger.Warn("Error fetching expired servers", err)
				continue
			}

			if len(servers) == 0 {
				continue
			}

			m.logger.Info("Fetched expired servers", servers)

			err = m.txManager.Do(ctx, func(ctx context.Context) error {
				for _, server := range servers {
					err = m.core.UnRentServer(ctx, server.ID)
					if err != nil {
						m.logger.Warn("Error un-renting server", err)
						return err
					}

					err = m.core.RebootServer(ctx, server.ID)
					if err != nil {
						m.logger.Warn("Error rebooting server", err)
						return err
					}
				}

				return nil
			})
			if err != nil {
				m.logger.Warn("Error rebooting servers", err)
			}
		case <-ctx.Done():
			m.logger.Info("Worker stopping")
		}
	}
}
