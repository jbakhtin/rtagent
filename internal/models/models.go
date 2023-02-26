package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/types"
	"strconv"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
}

type Valuer interface {
	Type() string
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

func (t *Metric) UnmarshalJSON(data []byte) error {
	// TODO: переделать в соответсвии с уроком: Спринт 2 -> Стандартные сериализаторы -> 2. Динамический JSON-объект
	// TODO: обраюотать ошибки если поля переданы неправильно
	var aliasValue Metrics

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	t.MKey = aliasValue.ID
	t.MType = aliasValue.MType

	switch t.MType {
	case types.GaugeType:
		t.MValue = *aliasValue.Value
	case types.CounterType:
		t.MValue = *aliasValue.Delta
	default:
		return errors.New("type not recognized")
	}

	return nil
}

func NewGauge(mType, mKey, mValue string) (Metricer, error){
	var Value types.Metricer

	switch mType {
	case types.GaugeType:
		floatV, err := strconv.ParseFloat(mValue, 64)
		if err != nil {
			return nil, err
		}
		Value = types.Gauge(floatV)
	case types.CounterType:
		intV, err := strconv.ParseInt(mValue, 10, 0)
		if err != nil {
			return nil, err
		}
		Value = types.Counter(intV)
	default:
		return nil, errors.New("types not recognized")
	}

	return Metric{
		MKey: mKey,
		MType: mType,
		MValue: Value,
	}, nil
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

func (g Gauge) Type() string {
	return g.Type()
}

func (c Counter) Type() string {
	return c.Type()
}

// Request
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *types.Counter   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *types.Gauge `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
