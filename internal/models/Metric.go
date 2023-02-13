package models

type Metric struct {
	Tp string
	K  string
	Vl string
}

func NewMetric(tp, k, vl string) Metric {
	return Metric{
		Tp: tp,
		K: k,
		Vl: vl,
	}
}

func (m Metric) Type() string {
	return m.Tp
}

func (m Metric) Key() string {
	return m.K
}

func (m Metric) Value() string {
	return m.Vl
}

