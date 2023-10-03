package server

import (
	"context"
	"crypto/rsa"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/pkg/crypto"
	"net/http"
	"net/http/pprof"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/server/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jbakhtin/rtagent/internal/server/handlers"
)

type MainServer struct {
	*http.Server
	repository storage.MetricRepository
}

func New(cfg config.Config, repository storage.MetricRepository) (MainServer, error) {
	server := &http.Server{
		Addr: cfg.Address,
	}

	return MainServer{
		Server: server,
		repository: repository,
	}, nil
}

func (ms MainServer) Start(ctx context.Context, cfg config.Config) (err error) {
	var handlerMetric *handlers.HandlerMetric
	r := chi.NewRouter()

	handlerMetric, err = handlers.NewHandlerMetric(ctx, cfg, ms.repository)
	if err != nil {
		return
	}

	var privateKey *rsa.PrivateKey

	if cfg.GetCryptoKey() != "" {
		privateKey, err = crypto.GetPrivateKey(cfg.GetCryptoKey())
		if err != nil {
			return
		}
	}

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.GZIPCompressor)
	//r.Use(middlewares.TrustedSubnet(cfg))

	r.Route("/", func(r chi.Router) {
		r.Get("/", handlerMetric.GetAllMetricsAsHTML())
		r.Get("/ping", handlerMetric.PingStorage())

		r.Route("/value/", func(r chi.Router) {
			r.Post("/", handlerMetric.GetMetricAsJSON()) // ToDo: перенести в отдельный пакет handlerMetricJSON
			r.Get("/{type}/{key}", handlerMetric.GetMetricValue())
		})

		r.Route("/update/", func(r chi.Router) {
			r.Use(middlewares.Decrypt(privateKey))
			r.Post("/", handlerMetric.UpdateMetricByJSON())
			r.Post("/{type}/{key}/{value}", handlerMetric.UpdateMetric())
		})

		r.Route("/updates/", func(r chi.Router) {
			r.Use(middlewares.Decrypt(privateKey))
			r.Post("/", handlerMetric.UpdateMetricsByJSON())
		})

		r.HandleFunc("/debug/pprof/*", pprof.Index)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	})

	go func() {
		if err := http.ListenAndServe(ms.Addr, r); err != nil {
			return
		}
	}()

	return
}
