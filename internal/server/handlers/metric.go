package handlers

import (
	"bytes"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/services"
	"html/template"
	"net/http"
)

type HandlerMetric struct {
	service *services.MetricService
}

var listOfMetricHTMLTemplate = `
	{{range .}}
			<div>{{.MKey}}:{{.StringValue}}</div>
	{{end}}
`

func NewHandlerMetric() (*HandlerMetric, error){
	service, err := services.NewMetricService()
	if err != nil {
		return nil, err
	}

	return &HandlerMetric{
		service: service,
	}, nil
}

func (h *HandlerMetric) GetGauge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}

		metric, err := h.service.GetGauge(mKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		mValue, err := metric.MarshalValue()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(mValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (h *HandlerMetric) GetCounter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}

		metric, err := h.service.GetCounter(mKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		mValue, err := metric.MarshalValue()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(mValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (h *HandlerMetric) UpdateCounter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mValue := chi.URLParam(r, "value")
		if mValue == "" {
			http.Error(w, "value not valid", http.StatusBadRequest)
			return
		}

		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "key not valid", http.StatusBadRequest)
			return
		}

		metric, err := models.NewCounter("counter", mKey, mValue)
		if err != nil {
			http.Error(w, "value not valid", http.StatusBadRequest)
			return
		}

		_, err = h.service.UpdateCounter(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) UpdateGauge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mValue := chi.URLParam(r, "value")
		if mValue == "" {
			http.Error(w, "value not valid", http.StatusBadRequest)
			return
		}

		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "key not valid", http.StatusBadRequest)
			return
		}

		metric, err := models.NewGauge("counter", mKey, mValue)
		if err != nil {
			http.Error(w, "value not valid", http.StatusBadRequest)
			return
		}

		_, err = h.service.UpdateGauge(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics, err := h.service.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		HTMLTemplate := template.Must(template.New("listOfMetricHTMLTemplate").Parse(listOfMetricHTMLTemplate))
		buffer := bytes.NewBuffer(nil)
		err = HTMLTemplate.Execute(buffer, metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = fmt.Fprint(w, buffer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
