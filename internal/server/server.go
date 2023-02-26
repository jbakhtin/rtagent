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
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAll())

		r.Route("/value/", func(r chi.Router) {
			r.Get("/", handlerMetric.GetV2()) //TODO: знать, стоит ли выносить хендлерами v2 в отдельный модуль
			r.Get("/{type}/{key}", handlerMetric.Get())
		})

		r.Route("/update/", func(r chi.Router) {
			r.Post("/", handlerMetric.UpdateV2()) // v2
			r.Post("/{type}/{key}/{value}", handlerMetric.Update())
		})
	})

	return http.ListenAndServe(s.serverAddress, r)
}
