package models

const  (
	GaugeType = "gauge"
	CounterType = "counter"
)

type (
	Gauge float64
	Counter int64

	Metric struct {
		MKey string
		MType string
		MGauge Gauge
		MCounter Counter
	}
)

func (g Gauge) Type() string {
	return GaugeType
}

func (c Counter) Type() string {
	return CounterType
}
