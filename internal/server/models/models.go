package models

import "github.com/jbakhtin/rtagent/internal/types"

type Metrics struct {
	MKey  string         `json:"id"`              // имя метрики
	MType string         `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func ToJSON(id string, value types.Metricer) Metrics {
	metric := Metrics{
		MKey:    id,
		MType: value.Type(),
	}
	switch v := value.(type) {
	case types.Counter:
		metric.Delta = &v
	case types.Gauge:
		metric.Value = &v
	}

	return metric
}




