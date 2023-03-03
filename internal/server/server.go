package server

import (
	"context"
	"net/http"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/server/handlers"
)

type Server struct {
	serverAddress string
}

func New(cfg config.Config) (Server, error) {
	return Server{
		serverAddress: cfg.Address,
	}, nil
}

func (s Server) Start(ctx context.Context, cfg config.Config) error {
	r := chi.NewRouter()

	handlerMetric, err := handlers.NewHandlerMetric(ctx, cfg)
	if err != nil {
		return err
	}

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	// TODO: узнать, можно ли реализовать через zap.Logger и как его лучше прокинть сюда, что бы логи были централизованы
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GZIPCompressor)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())

		r.Route("/value/", func(r chi.Router) {
			//TODO: узнать, стоит ли выносить хендлерами v2 в отдельный модуль, и как такое лучше делать
			r.Post("/", handlerMetric.GetV2())
			r.Get("/{type}/{key}", handlerMetric.Get())
		})

		r.Route("/update/", func(r chi.Router) {
			r.Post("/", handlerMetric.UpdateV2()) // v2
			r.Post("/{type}/{key}/{value}", handlerMetric.Update())
		})
	})

	return http.ListenAndServe(s.serverAddress, r)
}
