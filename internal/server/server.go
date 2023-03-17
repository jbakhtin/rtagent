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

type MainServer struct {
	*http.Server
}

func New(cfg config.Config) (MainServer, error) {
	server := &http.Server{
		Addr: cfg.Address,
	}

	return MainServer{
		Server: server,
	}, nil
}

func (ms MainServer) Start(ctx context.Context, cfg config.Config) error {
	r := chi.NewRouter()

	handlerMetric, err := handlers.NewHandlerMetric(ctx, cfg)
	if err != nil {
		return err
	}

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GZIPCompressor)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAllMetricsAsHTML())

		r.Route("/value/", func(r chi.Router) {
			r.Post("/", handlerMetric.GetMetricAsJSON()) // ToDo: перенести в отдельный пакет handlerMetricJSON
			r.Get("/{type}/{key}", handlerMetric.GetMetricValue())
		})

		r.Route("/update/", func(r chi.Router) {
			r.Post("/", handlerMetric.UpdateMetricByJSON())
			r.Post("/{type}/{key}/{value}", handlerMetric.UpdateMetric())
		})
	})

	return http.ListenAndServe(ms.Addr, r)
}
