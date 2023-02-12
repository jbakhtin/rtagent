package servers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/jbakhtin/rtagent/internal/servers/handlers"
	"net/http"
)

func Start() error {
	r := chi.NewRouter()

	ctx := context.Background()
	repo := inmemory.NewMetricRepository(&ctx)
	handlerMetric := handlers.NewHandlerMetric(repo)

	// middlewares
	r.Use(middleware.Logger)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.Get()) // TODO: Need return HTML

		r.Get("/value/{type}/{key}", handlerMetric.Find())
		r.Post("/update/{type}/{key}/{value}", handlerMetric.Update())
	})

	return http.ListenAndServe(":8080", r)
}
