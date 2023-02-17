package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/jbakhtin/rtagent/internal/server/handlers"
	"net/http"
)

func Start() error {
	r := chi.NewRouter()

	repo := inmemory.NewMetricRepository()
	handlerMetric := handlers.NewHandlerMetric(repo)

	// middlewares
	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())
		r.Get("/value/{type}/{key}", handlerMetric.Get())
		r.Post("/update/{type}/{key}/{value}", handlerMetric.Update())
	})

	return http.ListenAndServe(":8080", r)
}
