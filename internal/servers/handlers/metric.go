package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"github.com/jbakhtin/rtagent/internal/services"
	"net/http"
)

type HandlerMetric struct {
	repo interfaces.MetricRepository
}

func NewHandlerMetric(repo interfaces.MetricRepository) *HandlerMetric {
	return &HandlerMetric{
		repo: repo,
	}
}

func (h *HandlerMetric) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tp := chi.URLParam(r, "type")
		k := chi.URLParam(r, "key")

		metric, err := h.repo.Find(tp, k)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		metricJson, err := json.Marshal(metric.Value())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(metricJson)
	}
}

func (h *HandlerMetric) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tp := chi.URLParam(r, "type")
		if tp == "" || tp == "unknown" {
			http.Error(w, errors.New("type not valid").Error(), http.StatusNotImplemented)
			return
		}

		k := chi.URLParam(r, "key")
		// Need to add valudate

		vl := chi.URLParam(r, "value")
		if vl == "" || vl == "none" {
			http.Error(w, errors.New("value not valid").Error(), http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		service := services.NewMetricService(&ctx, h.repo)

		metric, err := service.Update(tp, k, vl)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metricJson, err := json.Marshal(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(metricJson)
	}
}

func (h *HandlerMetric) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := h.repo.Get()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		metricJson, err := json.Marshal(metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(metricJson)
	}
}