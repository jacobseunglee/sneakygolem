package logger

import (
	"log/slog"
	"os"
)

var (
	Server *slog.Logger
	Client *slog.Logger
)

func Init() {
	level := slog.LevelError

	Server = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).With("component", "server")
	Client = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})).With("component", "client")
}
