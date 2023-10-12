package models

import (
	"fmt"
	"github.com/jbakhtin/rtagent/pkg/hasher"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metrics struct {
	MKey  string         `json:"id"`
	MType string         `json:"type"`
	Delta *types.Counter `json:"delta,omitempty"`
	Value *types.Gauge   `json:"value,omitempty"`
	Hash  string         `json:"hash,omitempty"`
}

type Configer interface {
	GetKeyApp() string
}

func ToJSON(cfg Configer, id string, value types.Metricer) (Metrics, error) {
	var err error
	metric := Metrics{
		MKey:  id,
		MType: value.Type(),
	}
	switch v := value.(type) {
	case types.Counter:
		metric.Delta = &v
		metric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%x", metric.MKey, metric.MType, metric.Delta), []byte(cfg.GetKeyApp()))
	case types.Gauge:
		metric.Value = &v
		metric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%x", metric.MKey, metric.MType, metric.Value), []byte(cfg.GetKeyApp()))
	}
	if err != nil {
		return Metrics{}, err
	}

	return metric, nil
}
