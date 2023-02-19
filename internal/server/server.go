package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/jbakhtin/rtagent/internal/server/handlers"
)

func Start(serverAddress string) error {
	r := chi.NewRouter()

	repo := inmemory.NewMetricRepository()
	handlerMetric := handlers.NewHandlerMetric(repo)

	// middlewares
	r.Use(middleware.Logger) // TODO: need to add another middlewares

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())
		r.Get("/value/{type}/{key}", handlerMetric.Get())
		r.Post("/update/{type}/{key}/{value}", handlerMetric.Update())
	})

	return http.ListenAndServe(serverAddress, r)
}
