package types

const  (
	GaugeType = "gauge"
	CounterType = "counter"
)

type Metricer interface {
	Type() string
}

type (
	Gauge float64
	Counter int64
)

func (g Gauge) Type() string {
	return GaugeType
}

func (c Counter) Type() string {
	return CounterType
}
