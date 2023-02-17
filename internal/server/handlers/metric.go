package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"github.com/jbakhtin/rtagent/internal/services"
	"html/template"
	"net/http"
)

type HandlerMetric struct {
	repo interfaces.MetricRepository
}

var listOfMetricHTMLTemplate = `
	{{range .}}
			<div>{{ .K }}: {{ .Vl }}</div>
	{{end}}
`

func NewHandlerMetric(repo interfaces.MetricRepository) *HandlerMetric {
	return &HandlerMetric{
		repo: repo,
	}
}

func (h *HandlerMetric) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tp := chi.URLParam(r, "type")
		k := chi.URLParam(r, "key")

		metric, err := h.repo.Get(tp, k)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Write([]byte(metric.Value()))
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

		_, err := service.Update(tp, k, vl)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := h.repo.GetAll()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		t := template.Must(template.New("test").Parse(listOfMetricHTMLTemplate))
		buffer := bytes.NewBuffer(nil)
		err = t.Execute(buffer, metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, buffer)
	}
}