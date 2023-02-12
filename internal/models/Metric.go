package models

type Metric struct {
	tp string
	k  string
	vl string
}

func NewMetric(tp, k, vl string) Metric {
	return Metric{
		tp: tp,
		k: k,
		vl: vl,
	}
}

func (m Metric) Type() string {
	return m.tp
}

func (m Metric) Key() string {
	return m.k
}

func (m Metric) Value() string {
	return m.vl
}

