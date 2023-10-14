package models

import (
	"fmt"
	pb "github.com/jbakhtin/rtagent/gen/go/metric/v1"
	"github.com/jbakhtin/rtagent/pkg/hasher"
	"strconv"

	"github.com/jbakhtin/rtagent/internal/server/models"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Metricer interface {
	Type() string
	Key() string
	StringValue() string
	ToJSON(key []byte) (models.Metrics, error)
}

type (
	Description struct {
		MKey  string
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

func NewGauge(mType, mKey, mValue string) (Gauge, error) {
	value, err := strconv.ParseFloat(mValue, 64)
	if err != nil {
		return Gauge{}, err
	}

	return Gauge{
		Description: Description{
			MKey:  mKey,
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

func (g Gauge) ToJSON(key []byte) (models.Metrics, error) {
	var err error
	JSONMetric := models.Metrics{
		MKey:  g.MKey,
		MType: g.MType,
		Value: &g.MValue,
	}
	if len(key) != 0 {
		JSONMetric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%f", JSONMetric.MKey, JSONMetric.MType, *JSONMetric.Value), key)
		if err != nil {
			return models.Metrics{}, err
		}
	}

	return JSONMetric, nil
}

func (g Gauge) ToGRPC(key string) (*pb.Metric, error) {
	var err error
	GRPCMetric := &pb.Metric{
		Key:   g.MKey,
		Type:  pb.Type_TYPE_GAUGE,
		Value: float32(g.MValue),
	}

	if len(key) != 0 {
		GRPCMetric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%f", GRPCMetric.Key, GRPCMetric.Type, GRPCMetric.Value), []byte(key))
		if err != nil {
			return nil, err
		}
	}

	return GRPCMetric, nil
}

// Counter ----

func NewCounter(mType, mKey, mValue string) (Counter, error) {
	value, err := strconv.ParseInt(mValue, 10, 0)
	if err != nil {
		return Counter{}, err
	}

	return Counter{
		Description: Description{
			MKey:  mKey,
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
	value := fmt.Sprintf("%d", c.MValue)
	return value
}

func (c Counter) ToJSON(key []byte) (models.Metrics, error) {
	var err error
	JSONMetric := models.Metrics{
		MKey:  c.MKey,
		MType: c.MType,
		Delta: &c.MValue,
	}

	if len(key) != 0 {
		JSONMetric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%d", JSONMetric.MKey, JSONMetric.MType, *JSONMetric.Delta), key)
		if err != nil {
			return models.Metrics{}, err
		}
	}

	return JSONMetric, nil
}

func (c Counter) ToGRPC(key string) (*pb.Metric, error) {
	var err error
	GRPCMetric := &pb.Metric{
		Key:   c.MKey,
		Type:  pb.Type_TYPE_COUNTER,
		Delta: uint64(c.MValue),
	}

	if len(key) != 0 {
		GRPCMetric.Hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%d", GRPCMetric.Key, GRPCMetric.Type, GRPCMetric.Delta), []byte(key))
		if err != nil {
			return nil, err
		}
	}

	return GRPCMetric, nil
}

func (c *Counter) Increment() {
	c.MValue++
}

func (c *Counter) Add(value types.Counter) {
	c.MValue += value
}
