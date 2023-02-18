package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"github.com/jbakhtin/rtagent/internal/services"
	"html/template"
	"net/http"
	"strconv"
)

type gauge float64
type counter int64

type HandlerMetric struct {
	repo interfaces.MetricRepository
}

var listOfMetricHTMLTemplate = `
	{{range .}}
		{{if (eq .MType "gauge")}}
			<div>{{.MKey}}:{{.MGauge}}</div>
		{{else if (eq .MType "counter")}}
			<div>{{.MKey}}:{{.MCounter}}</div>
		{{end}}
	{{end}}
`

func NewHandlerMetric(repo interfaces.MetricRepository) *HandlerMetric {
	return &HandlerMetric{
		repo: repo,
	}
}

func (h *HandlerMetric) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var mValue []byte
		var metric models.Metric

		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "record not found", http.StatusNotFound)
		}

		metric, err = h.repo.Get(mKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		switch metric.MType {
			case "gauge":
				mValue, err = json.Marshal(metric.MGauge)
			case "counter":
				mValue, err = json.Marshal(metric.MCounter)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = w.Write(mValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (h *HandlerMetric) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := models.Metric{}

		mValue := chi.URLParam(r, "value")
		if mValue == "" {
			http.Error(w, errors.New("value not valid").Error(), http.StatusBadRequest)
			return
		}

		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, errors.New("value not valid").Error(), http.StatusBadRequest)
			return
		}

		mType := chi.URLParam(r, "type")
		switch mType {
			case "gauge":
				floatValue, err := strconv.ParseFloat(mValue, 64)
				if err != nil {
					http.Error(w, errors.New("value not valid").Error(), http.StatusBadRequest)
					return
				}

				metric.MGauge = models.Gauge(floatValue)
				metric.MType = metric.MGauge.Type()
			case "counter":
				intValue, err := strconv.ParseInt(mValue, 10, 0)
				if err != nil {
					http.Error(w, errors.New("value not valid").Error(), http.StatusBadRequest)
					return
				}

				metric.MCounter = models.Counter(intValue)
				metric.MType = metric.MCounter.Type()
		default:
			http.Error(w, errors.New("type not valid").Error(), http.StatusNotImplemented)
			return
		}

		metric.MKey = mKey

		ctx := context.Background()
		service := services.NewMetricService(&ctx, h.repo)

		_, err := service.Update(metric)
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

		fmt.Println(metrics)

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