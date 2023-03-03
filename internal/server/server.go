package server

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server/middlewares"
	"net/http"

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
	r.Use(middleware.Logger) // TODO: узнать, можно ли реализовать через zap.Logger
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GZIPCompress)
	r.Use(middlewares.GZIPCompress)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())

		r.Route("/value/", func(r chi.Router) {
			r.Post("/", handlerMetric.GetV2()) //TODO: узнать, стоит ли выносить хендлерами v2 в отдельный модуль
			r.Get("/{type}/{key}", handlerMetric.Get())
		})

		r.Route("/update/", func(r chi.Router) {
			r.Post("/", handlerMetric.UpdateV2()) // v2
			r.Post("/{type}/{key}/{value}", handlerMetric.Update())
		})
	})

	fmt.Println(s.serverAddress)

	return http.ListenAndServe(s.serverAddress, r)
}
