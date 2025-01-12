package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/lulzshadowwalker/recall/internal/http/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		slog.Error("app creation failed", "err", err)
		return
	}

	slog.Info("server started", "addr", application.Addr(), "read_timeout", application.ReadTimeout(), "write_timeout", application.WriteTimeout())
	if err := application.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server shutdown", "err", err)
			return
		}

		slog.Info("server crashed", "err", err)
	}
}
