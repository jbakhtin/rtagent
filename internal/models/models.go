package models

import (
	"fmt"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
}

type Metric struct {
	MKey    string       `json:"id"`              // имя метрики
	MType string         `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Metric) Type() string {
	return m.MType
}

func (m Metric) Key() string {
	return m.MKey
}

func (m Metric) StringValue() string {
	switch m.MType {
	case types.CounterType:
		return fmt.Sprintf("%v", *m.Delta)
	case types.GaugeType:
		return fmt.Sprintf("%v", *m.Value)
	}

	return ""
}
