package app

import (
	"errors"
	"net/http"
	"regexp"
	"time"

	"github.com/lulzshadowwalker/recall/internal/http/handler"
)

const (
	AppDefaultReadTimeout  time.Duration = 2 * time.Second
	AppDefaultWriteTimeout time.Duration = 2 * time.Second
	AppDefaultAddr         string        = ":8080"
)

type App struct {
	router *http.ServeMux
	server *http.Server
}

type AppOption func(*App) error

// TODO: This does not belong here
// would be nice to have something like homeHandler.Index.Unwrap()
type DecoratedHandler func(w http.ResponseWriter, r *http.Request) error

func New(opts ...AppOption) (*App, error) {
	router := http.NewServeMux()
	homeHandler := handler.HomeHandler{}
	router.HandleFunc("/", homeHandler.Index) //  TODO: Add decorated handler

	app := &App{
		router: router,
		server: &http.Server{
			Addr:         AppDefaultAddr,
			Handler:      router,
			ReadTimeout:  AppDefaultReadTimeout,
			WriteTimeout: AppDefaultWriteTimeout,
		},
	}

	for _, opt := range opts {
		if err := opt(app); err != nil {
			return nil, err
		}
	}

	return app, nil
}

func (a *App) Start() error {
	return a.server.ListenAndServe()
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

		a.server.Addr = addr
		return nil
	}
}

func WithReadTimeout(d time.Duration) AppOption {
	return func(a *App) error {
		a.server.ReadTimeout = d
		return nil
	}
}

func (a *App) WithWriteTimeout(d time.Duration) AppOption {
	return func(a *App) error {
		a.server.WriteTimeout = d
		return nil
	}
}

func (a *App) Addr() string {
	return a.server.Addr
}

func (a *App) ReadTimeout() time.Duration {
	return a.server.ReadTimeout
}

func (a *App) WriteTimeout() time.Duration {
	return a.server.WriteTimeout
}
