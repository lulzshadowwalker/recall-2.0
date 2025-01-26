package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lulzshadowwalker/recall/internal/http/handler"
	"github.com/lulzshadowwalker/recall/internal/psql"
)

const (
	AppDefaultReadTimeout  time.Duration = 2 * time.Second
	AppDefaultWriteTimeout time.Duration = 2 * time.Second
	AppDefaultAddr         string        = ":8080"
)

type App struct {
	Echo    *echo.Echo
	addr    string
	timeout time.Duration
	db      *pgxpool.Pool
}

type AppOption func(*App) error

func New(opts ...AppOption) (*App, error) {
	e := echo.New()

	app := &App{
		timeout: AppDefaultReadTimeout,
		addr:    AppDefaultAddr,
		Echo:    e,
	}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	if app.db == nil {
		return nil, errors.New("db cannot be nil")
	}

	{
		mh := handler.NewMemoryHandler(psql.NewMemory(app.db))
		e.POST("/memories", mh.Create)
		e.GET("/memories", mh.Index)
		e.DELETE("/memories/:memory", mh.Delete)
	}

	//  NOTE: Middlewares should be added after all options are applied
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: app.timeout,
	}))

	e.Validator = NewRecallValidator()

	e.HTTPErrorHandler = recallHTTPErrorHandler

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogError:  true,
		// HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	return app, nil
}

func (a *App) Start() error {
	return a.Echo.Start(a.addr)
}

func (a *App) WithAddr(addr string) AppOption {
	return func(a *App) error {
		if addr == "" {
			return errors.New("addr cannot be empty")
		}

		regex := `^(:\d{1,5})$`
		if !regexp.MustCompile(regex).MatchString(addr) {
			return errors.New("addr must be in format :<port>")
		}

		a.addr = addr
		return nil
	}
}

func WithTimeout(d time.Duration) AppOption {
	return func(a *App) error {
		if d < 0 {
			return errors.New("timeout cannot be negative")
		}

		a.timeout = d

		return nil
	}
}

func WithDB(db *pgxpool.Pool) AppOption {
	return func(a *App) error {
		a.db = db
		return nil
	}
}

func (a *App) Close() {
	a.db.Close()
}

func (a *App) Addr() string {
	return a.addr
}

func (a *App) Timeout() time.Duration {
	return a.timeout
}
