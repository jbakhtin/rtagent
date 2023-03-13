package models

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
}

type Metric struct {
	MKey  string         `json:"id"`              // имя метрики
	MType string         `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string         `json:"hash,omitempty"`  //  значение хэша от MKey:MType:Delta|Value
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

func (m Metric) CalcHash(key []byte) (string, error) {
	h := hmac.New(sha256.New, key)

	switch m.MType {
	case types.CounterType:
		h.Write([]byte(fmt.Sprintf("%s:%s:%v", m.MKey, m.MType, &m.Delta)))
	case types.GaugeType:
		h.Write([]byte(fmt.Sprintf("%s:%s:%v", m.MKey, m.MType, &m.Value)))
	}

	dst := h.Sum(nil)
	return fmt.Sprintf("%x", dst), nil

}
