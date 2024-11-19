package main

import (
	"log/slog"
	"machineIssuerSystem/internal/app"
	logTools "machineIssuerSystem/pkg/logger"
	"os"
)

func main() {
	logger := logTools.NewLogger(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	application, err := app.NewProductorApp(logger)
	if err != nil {
		logger.Error("App don't created", slog.Any("error", err))
	}
	err = application.Start()
	if err != nil {
		logger.Error("App don't started", slog.Any("error", err))
	}
}
