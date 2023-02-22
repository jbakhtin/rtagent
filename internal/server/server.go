package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/server/handlers"
)

type Server struct {
	serverAddress string
}

func New(serverAddress string) (Server, error) {
	return Server{
		serverAddress: serverAddress,
	}, nil
}

func (s Server) Start() error {
	r := chi.NewRouter()

	handlerMetric, err := handlers.NewHandlerMetric()
	if err != nil {
		return err
	}

	// middlewares
	r.Use(middleware.Logger) // TODO: need to add another middlewares

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())
		r.Get("/value/gauge/{key}", handlerMetric.GetGauge())
		r.Get("/value/counter/{key}", handlerMetric.GetCounter())
		r.Post("/update/gauge/{key}/{value}", handlerMetric.UpdateGauge())
		r.Post("/update/counter/{key}/{value}", handlerMetric.UpdateCounter())
	})

	return http.ListenAndServe(s.serverAddress, r)
}
