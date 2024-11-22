package main

import (
	"log/slog"
	"os"

	"machineIssuerSystem/internal/app"
	logTools "machineIssuerSystem/pkg/logger"
)

func main() {
	logger := logTools.NewLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	application, err := app.NewApplication(logger)
	if err != nil {
		logger.Error("App don't created", slog.Any("error", err))
	}
	err = application.Start()
	if err != nil {
		logger.Error("App don't started", slog.Any("error", err))
	}
}
