package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/pkg/hasher"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	handlerModels "github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
)

// MetricRepository интерфейс реализации хранилища.
type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Set(models.Metricer) (models.Metricer, error)
	SetBatch([]models.Metricer) ([]models.Metricer, error)
	TestPing() error
}

type HandlerMetric struct {
	repository MetricRepository
	config     config.Config
}

// listOfMetricHTMLTemplate html шаблон списка метрик key: value.
var listOfMetricHTMLTemplate = `
	{{range .}}
		<div>{{.Key}}:{{.StringValue}}</div>
	{{end}}
`

func NewHandlerMetric(ctx context.Context, cfg config.Config, repository MetricRepository) (*HandlerMetric, error) {
	return &HandlerMetric{
		repository: repository,
		config:     cfg,
	}, nil
}

// GetMetricValue - получить значение метрики по типу метрики и ключу.
//
//	/value/{type}/{key}
//	type - тип метрики
//	key - ключ метрики
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

		metric, err := h.repository.Get(mKey)
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

// GetMetricAsJSON - получить значение метрики по типу метрики и ключу в формате JSON.
//
//	/value/{type}/{key}
//	type - тип метрики
//	key - ключ метрики
func (h *HandlerMetric) GetMetricAsJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var metrics handlerModels.Metrics
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		metric, err := h.repository.Get(metrics.MKey)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var JSONMetric handlerModels.Metrics
		JSONMetric, err = metric.ToJSON([]byte(h.config.KeyApp))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var jsonMetric []byte
		jsonMetric, err = json.Marshal(JSONMetric)
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

// UpdateMetric - обновить значение метрики по типу метрики и ключу в формате JSON.
//
//	/value/{type}/{key}/{value}
//	type - тип метрики
//	key - ключ метрики
//	value - новое значение метрики
func (h *HandlerMetric) UpdateMetric() http.HandlerFunc {
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

		var metric models.Metricer
		var err error

		mType := chi.URLParam(r, "type")
		switch mType {
		case types.GaugeType:
			metric, err = models.NewGauge(mType, mKey, mValue)
		case types.CounterType:
			metric, err = models.NewCounter(mType, mKey, mValue)
		default:
			http.Error(w, "type not valid", http.StatusNotImplemented)
			return
		}
		if err != nil {
			http.Error(w, "type not valid", http.StatusBadRequest)
			return
		}

		_, err = h.repository.Set(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// UpdateMetricByJSON - создать или обновить несколько метрик в формате JSON
//
//	/value/{type}/{key}/{value}
//	type - тип метрики
//	key - ключ метрики
//	value - новое значение метрики
func (h *HandlerMetric) UpdateMetricByJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var metrics handlerModels.Metrics
		var hash string

		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		var metric models.Metricer
		switch metrics.MType {
		case types.GaugeType:
			metric, err = models.NewGauge(metrics.MType, metrics.MKey, fmt.Sprintf("%v", *metrics.Value))
			if h.config.KeyApp != "" {
				hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%f", metrics.MKey, metrics.MType, *metrics.Value), []byte(h.config.KeyApp))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if hash != metrics.Hash {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		case types.CounterType:
			metric, err = models.NewCounter(metrics.MType, metrics.MKey, fmt.Sprintf("%v", *metrics.Delta))
			if h.config.KeyApp != "" {
				hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%d", metrics.MKey, metrics.MType, *metrics.Delta), []byte(h.config.KeyApp))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if hash != metrics.Hash {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = h.repository.Set(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		JSONMetric, err := metric.ToJSON([]byte(h.config.KeyApp))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonMetric, err := json.Marshal(JSONMetric)
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
		metrics, err := h.repository.GetAll()
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

func (h *HandlerMetric) PingStorage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.repository.TestPing()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) UpdateMetricsByJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var metrics []handlerModels.Metrics
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		var hash string
		mMetrics := make([]models.Metricer, len(metrics))
		for i, m := range metrics {
			var metric models.Metricer
			switch m.MType {
			case types.GaugeType:
				metric, err = models.NewGauge(m.MType, m.MKey, fmt.Sprintf("%v", *m.Value))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if h.config.KeyApp != "" {
					hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%f", m.MKey, m.MType, *m.Value), []byte(h.config.KeyApp))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if hash != m.Hash {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
			case types.CounterType:
				metric, err = models.NewCounter(m.MType, m.MKey, fmt.Sprintf("%d", *m.Delta))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if h.config.KeyApp != "" {
					hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%d", m.MKey, m.MType, *m.Delta), []byte(h.config.KeyApp))
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					if hash != m.Hash {
						w.WriteHeader(http.StatusBadRequest)
						return
					}
				}
			}

			mMetrics[i] = metric
		}

		_, err = h.repository.SetBatch(mMetrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte("[]"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
