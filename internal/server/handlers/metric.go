package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"github.com/jbakhtin/rtagent/internal/services"
	"html/template"
	"net/http"
	"strconv"
)

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
			case models.GaugeType:
				mValue, err = json.Marshal(metric.MGauge)
			case models.CounterType:
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
			http.Error(w,"value not valid", http.StatusBadRequest)
			return
		}

		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w,"key not valid", http.StatusBadRequest)
			return
		}

		mType := chi.URLParam(r, "type")
		switch mType {
			case models.GaugeType:
				floatValue, err := strconv.ParseFloat(mValue, 64)
				if err != nil {
					http.Error(w, "value not valid", http.StatusBadRequest)
					return
				}

				metric.MGauge = models.Gauge(floatValue)
				metric.MType = metric.MGauge.Type()
			case models.CounterType:
				intValue, err := strconv.ParseInt(mValue, 10, 0)
				if err != nil {
					http.Error(w,"value not valid", http.StatusBadRequest)
					return
				}

				metric.MCounter = models.Counter(intValue)
				metric.MType = metric.MCounter.Type()
		default:
			http.Error(w, "type not valid", http.StatusNotImplemented)
			return
		}

		metric.MKey = mKey

		service := services.NewMetricService(h.repo)

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