package xlog

import (
	"io"
	"log/slog"
)

var options = &slog.HandlerOptions{Level: slog.LevelDebug}

func New(writer io.Writer, fields Fields) *slog.Logger {
	return slog.New(slog.NewJSONHandler(writer, options).WithAttrs(fields.ToAttrs()))
}
