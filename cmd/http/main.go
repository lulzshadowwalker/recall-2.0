package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	application "github.com/lulzshadowwalker/recall/internal/http/app"
	"github.com/lulzshadowwalker/recall/internal/psql"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		slog.Error("failed to load .env.local", "err", err)
		os.Exit(1)
	}

	pool, err := psql.Connect(psql.ConnectionParams{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	app, err := application.New(application.WithDB(pool))
	if err != nil {
		slog.Error("app creation failed", "err", err)
		return
	}

	app.Echo.Logger.Info("server started", "addr", app.Addr(), "timeout", app.Timeout())
	if err := app.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server shutdown", "err", err)
		} else {
			slog.Info("server crashed", "err", err)
		}
	}
	//  TODO: Graceful termination
	defer app.Close()
}
