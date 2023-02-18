package models

type Gauge float64
type Counter int64

type Metric struct {
	MKey string
	MType string
	MGauge Gauge
	MCounter Counter
}

func (g Gauge) Type() string {
	return "gauge"
}

func (c Counter) Type() string {
	return "counter"
}
