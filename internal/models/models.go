package models

import (
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/types"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
}

type (
	Metric struct {
		MKey   string      `json:"id,omitempty"`
		MType  string      `json:"type,omitempty"`
		MValue interface{}
	}

	Gauge struct {
		VValue types.Gauge `json:"value"`
	}

	Counter struct {
		VValue types.Counter `json:"delta"`
	}
)

func (m *Metric) UnmarshalJSON(data []byte) error {
	// TODO: переделать в соответсвии с уроком: Спринт 2 -> Стандартные сериализаторы -> 2. Динамический JSON-объект
	// TODO: обраюотать ошибки если поля переданы неправильно
	var aliasValue Metrics

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	fmt.Println(aliasValue)

	m.MKey = aliasValue.ID
	m.MType = aliasValue.MType

	switch m.MType {
	case types.GaugeType:
		if aliasValue.Value != nil  {
			m.MValue = *aliasValue.Value
		}
	case types.CounterType:
		if aliasValue.Delta != nil  {
			m.MValue = *aliasValue.Delta
		}
	}

	return nil
}

func (m Metric) Type() string {
	return m.MType
}

func (m Metric) Key() string {
	return m.MKey
}

func (m Metric) StringValue() string {
	value := fmt.Sprintf("%v", m.MValue)
	return value
}

// Request
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
