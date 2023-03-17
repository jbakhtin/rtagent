package models

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/server/models"
	"strconv"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
	ToJSON() models.Metrics
}

type (
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
	return g.MType
}

func (g Gauge) Key() string {
	return g.MKey
}

func (g Gauge) StringValue() string {
	return fmt.Sprintf("%v", g.MValue)
}

func (g Gauge) ToJSON() models.Metrics {
	return models.Metrics{
		MKey:  g.MKey,
		MType: g.MType,
		Value: &g.MValue,
	}
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
	return c.MType
}

func (c Counter) Key() string {
	return c.MKey
}

func (c Counter) StringValue() string {
	value := fmt.Sprintf("%v", c.MValue)
	return value
}

func (c Counter) ToJSON() models.Metrics {
	return models.Metrics{
		MKey:  c.MKey,
		MType: c.MType,
		Delta: &c.MValue,
	}
}

func (c *Counter) Increment() {
	c.MValue++
}

func (c *Counter) Add(value types.Counter) {
	c.MValue += value
}
