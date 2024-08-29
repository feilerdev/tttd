package main

import "log/slog"

var logger *slog.Logger

func NewLogger() {
	logger = slog.Default()
}
