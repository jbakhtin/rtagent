package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/services"
	"github.com/jbakhtin/rtagent/internal/types"
	"html/template"
	"net/http"
	"strconv"
)

type HandlerMetric struct {
	service *services.MetricService
}

var listOfMetricHTMLTemplate = `
	{{range .}}
		<div>{{.Key}}:{{.StringValue}}</div>
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

func (h *HandlerMetric) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mKey := chi.URLParam(r, "key")
		if mKey == "" {
			http.Error(w, "record not found", http.StatusNotFound)
			return
		}

		mType := chi.URLParam(r, "type")
		if mType == "" {
			http.Error(w, "invalid type", http.StatusInternalServerError)
			return
		}

		metric, err := h.service.Get(mKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		_, err = w.Write([]byte(metric.StringValue()))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}

func (h *HandlerMetric) GetV2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var metrics models.Metric
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metric, err := h.service.Get(metrics.MKey)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var buf bytes.Buffer
		jsonEncoder := json.NewEncoder(&buf)
		err = jsonEncoder.Encode(metric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err = w.Write(buf.Bytes())
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (h *HandlerMetric) Update() http.HandlerFunc {
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

		mType := chi.URLParam(r, "type")

		var Value interface{}
		switch mType {
		case types.GaugeType:
			floatV, err := strconv.ParseFloat(mValue, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			Value = types.Gauge(floatV)
		case types.CounterType:
			intV, err := strconv.ParseInt(mValue, 10, 0)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			Value = types.Counter(intV)
		default:
			http.Error(w, "type not valid", http.StatusNotImplemented)
			return
		}

		metric := models.Metric{
			MKey: mKey,
			MType: mType,
			MValue: Value,
		}

		_, err := h.service.Update(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) UpdateV2() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metric
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if _, err := h.service.Update(metric); err != nil {
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
