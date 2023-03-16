package models

import "github.com/jbakhtin/rtagent/internal/types"

type Metrics struct {
	MKey  string         `json:"id"`              // имя метрики
	MType string         `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}





