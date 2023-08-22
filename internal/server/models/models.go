package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/jbakhtin/rtagent/internal/config"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metrics struct {
	MKey  string         `json:"id"`
	MType string         `json:"type"`
	Delta *types.Counter `json:"delta,omitempty"`
	Value *types.Gauge   `json:"value,omitempty"`
	Hash  string         `json:"hash,omitempty"`
}

func ToJSON(cfg config.Config, id string, value types.Metricer) (Metrics, error) {
	var err error
	metric := Metrics{
		MKey:  id,
		MType: value.Type(),
	}
	switch v := value.(type) {
	case types.Counter:
		metric.Delta = &v
		metric.Hash, err = metric.CalcHash([]byte(cfg.KeyApp))
	case types.Gauge:
		metric.Value = &v
		metric.Hash, err = metric.CalcHash([]byte(cfg.KeyApp))
	}
	if err != nil {
		return Metrics{}, err
	}

	return metric, nil
}

func (m Metrics) CalcHash(key []byte) (string, error) {
	h := hmac.New(sha256.New, key)

	switch m.MType {
	case types.CounterType:
		h.Write([]byte(fmt.Sprintf("%s:%s:%d", m.MKey, m.MType, *m.Delta)))
	case types.GaugeType:
		h.Write([]byte(fmt.Sprintf("%s:%s:%f", m.MKey, m.MType, *m.Value)))
	}

	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst), nil
}
