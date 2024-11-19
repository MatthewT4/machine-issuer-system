package logger

import (
	"io"
	"log/slog"
)

func NewLogger(w io.Writer, opts *slog.HandlerOptions) *slog.Logger {
	h := ContextHandler{slog.NewJSONHandler(w, opts)}
	return slog.New(h)
}
