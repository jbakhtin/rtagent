package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/services"
	"github.com/jbakhtin/rtagent/internal/types"
)

type HandlerMetric struct {
	service *services.MetricService
}

var listOfMetricHTMLTemplate = `
	{{range .}}
		<div>{{.Key}}:{{.StringValue}}</div>
	{{end}}
`

func NewHandlerMetric(ctx context.Context, cfg config.Config) (*HandlerMetric, error) {
	service, err := services.NewMetricService(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &HandlerMetric{
		service: service,
	}, nil
}

func (h *HandlerMetric) GetMetricValue() http.HandlerFunc {
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

func (h *HandlerMetric) GetMetricAsJSON() http.HandlerFunc {
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

func (h *HandlerMetric) UpdateMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mValue := chi.URLParam(r, "value")
		if mValue == "" {
			// TODO: вынести ошибки в констаны
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
			MKey:   mKey,
			MType:  mType,
			MValue: Value,
		}

		test, err := h.service.Update(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonMetric, err := json.Marshal(test)
		fmt.Println(jsonMetric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err = w.Write(jsonMetric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func (h *HandlerMetric) UpdateMetricByJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var metric models.Metric
		err := json.NewDecoder(r.Body).Decode(&metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		test, err := h.service.Update(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonMetric, err := json.Marshal(test)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err = w.Write(jsonMetric)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) GetAllMetricsAsHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
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
