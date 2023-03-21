package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	models2 "github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/storages/filestorage"
	"github.com/jbakhtin/rtagent/internal/storages/postgres"
	"github.com/jbakhtin/rtagent/internal/types"
	"html/template"
	"net/http"
)

type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Set(models.Metricer) (models.Metricer, error)
	SetBatch(context.Context, []models.Metricer) ([]models.Metricer, error)
}

type HandlerMetric struct {
	repository MetricRepository
}

var listOfMetricHTMLTemplate = `
	{{range .}}
		<div>{{.Key}}:{{.StringValue}}</div>
	{{end}}
`

func NewHandlerMetric(ctx context.Context, cfg config.Config) (*HandlerMetric, error) {
	if cfg.DatabaseDSN != "" {
		ms, err := postgres.New("pgx", cfg) // TODO: move to cfg
		if err != nil {
			return nil, err
		}

		return &HandlerMetric{
			repository: &ms,
		}, nil
	}

	ms, err := filestorage.New(cfg)
	if err != nil {
		return nil, err
	}

	err = ms.Start(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &HandlerMetric{
		repository: &ms,
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

func (h *HandlerMetric) GetMetricAsJSON(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var metrics models2.Metrics
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

		JSONMetric, _ := metric.ToJSON([]byte(cfg.KeyApp))
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
	}
}

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

func (h *HandlerMetric) UpdateMetricByJSON(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var metrics models2.Metrics
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		if cfg.KeyApp != "" {
			hash, err := metrics.CalcHash([]byte(cfg.KeyApp))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if hash != metrics.Hash {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var metric models.Metricer
		switch metrics.MType {
		case types.GaugeType:
			metric, err = models.NewGauge(metrics.MType, metrics.MKey, fmt.Sprintf("%v", *metrics.Value))
		case types.CounterType:
			metric, err = models.NewCounter(metrics.MType, metrics.MKey, fmt.Sprintf("%v", *metrics.Delta))
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		test, err := h.repository.Set(metric)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		JSONMetric, err := test.ToJSON([]byte(cfg.KeyApp))
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

func (h *HandlerMetric) TestDBConnection(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := pgx.Connect(context.Background(), cfg.DatabaseDSN)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer conn.Close(context.Background())

		var id string
		var value string
		err = conn.QueryRow(context.Background(), "select id, value from metrics where id=$1", "Alloc").Scan(&id, &value)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (h *HandlerMetric) UpdateMetricsByJSON(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var metrics []models2.Metrics
		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}

		mMetrics := make([]models.Metricer, len(metrics))
		for i, m := range metrics {
			if cfg.KeyApp != "" {
				hash, err := m.CalcHash([]byte(cfg.KeyApp))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				if hash != m.Hash {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
			}

			var metric models.Metricer
			switch m.MType {
			case types.GaugeType:
				metric, err = models.NewGauge(m.MType, m.MKey, fmt.Sprintf("%v", *m.Value))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			case types.CounterType:
				metric, err = models.NewCounter(m.MType, m.MKey, fmt.Sprintf("%v", *m.Delta))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			mMetrics[i] = metric
		}

		_, err = h.repository.SetBatch(context.TODO(), mMetrics)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte("{}"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
