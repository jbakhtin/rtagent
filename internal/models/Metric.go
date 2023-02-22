package models

import (
	"encoding/json"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/types"
	"strconv"
)

type Metricer interface {
 	MarshalValue() ([]byte, error)
	StringValue() string
}

type Description struct {
	MKey string
	MType string
}

type Gauge struct {
	Description
	MValue types.Gauge
}

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

func (g Gauge) MarshalValue() ([]byte, error)  {
	bytes, err := json.Marshal(g.MValue)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (g Gauge) StringValue() string  {
	return fmt.Sprintf("%v", g.MValue)
}

type Counter struct {
	Description
	MValue types.Counter
}

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

func (c Counter) MarshalValue() ([]byte, error)  {
	bytes, err := json.Marshal(c.MValue)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (c Counter) StringValue() string  {
	return fmt.Sprintf("%v", c.MValue)
}

func (c *Counter) Increment()  {
	c.MValue++
}

func (c *Counter) Flush()  {
	c.MValue = 0
}

func (c *Counter) Add(counter types.Counter)  {
	c.MValue += counter
}