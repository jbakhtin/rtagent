package models

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/types"
	"strconv"
)

const  (
	GaugeType = "gauge"
	CounterType = "counter"
)

type (
	Metricer interface {
		Type() string
		Key() string
		StringValue() string
	}

	Description struct {
		MKey string
		MType string
	}

	Gauge struct {
		Description
		MValue types.Gauge
	}

	Counter struct {
		Description
		MValue types.Counter
	}
)

// Gauge ----

func NewGauge(mType, mKey, mValue string) (Gauge, error){
	value, err := strconv.ParseFloat(mValue, 64)
	if err != nil {
		return Gauge{}, err
	}

	return Gauge{
		Description: Description {
			MKey: mKey,
			MType: mType,
		},
		MValue: types.Gauge(value),
	}, nil
}


func (g Gauge) Type() string {
	return GaugeType
}

func (g Gauge) Key() string {
	return g.MKey
}

func (g Gauge) StringValue() string {
	value := fmt.Sprintf("%v", g.MValue)
	return value
}

// Counter ----

func NewCounter(mType, mKey, mValue string) (Counter, error){
	value, err := strconv.ParseInt(mValue, 10, 0)
	if err != nil {
		return Counter{}, err
	}

	return Counter{
		Description: Description {
			MKey: mKey,
			MType: mType,
		},
		MValue: types.Counter(value),
	}, nil
}

func (c Counter) Type() string {
	return CounterType
}

func (c Counter) Key() string {
	return c.MKey
}

func (c Counter) StringValue() string {
	value := fmt.Sprintf("%v", c.MValue)
	return value
}

func (c *Counter) Increment() {
	c.MValue++
}

func (c *Counter) Add(value types.Counter) {
	c.MValue += value
}
