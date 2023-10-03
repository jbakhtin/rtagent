package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

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
		metric.Hash, err = metric.CalcHash([]byte(cfg.GetKeyApp()))
	case types.Gauge:
		metric.Value = &v
		metric.Hash, err = metric.CalcHash([]byte(cfg.GetKeyApp()))
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
		_, err := h.Write([]byte(fmt.Sprintf("%s:%s:%d", m.MKey, m.MType, *m.Delta)))
		if err != nil {
			return "", err
		}
	case types.GaugeType:
		_, err := h.Write([]byte(fmt.Sprintf("%s:%s:%f", m.MKey, m.MType, *m.Value)))
		if err != nil {
			return "", err
		}
	}

	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst), nil
}
