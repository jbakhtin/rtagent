package aggregator

import "github.com/jbakhtin/rtagent/internal/types"

type Builder struct {
	aggregator aggregator
	err error
}

func New() *Builder {
	return &Builder{
		aggregator: aggregator{
			collection: Metrics{
				items: make(map[string]types.Metricer, 0),
			},
		},
	}
}

func (b *Builder) WithDefaultCollectors() *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, []CollectorFunc{Runtime, Gopsutil, RandomMetric, b.aggregator.poolCountCollector}...)
	return b
}

func (b *Builder) WithCustomCollector(collector CollectorFunc) *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, collector)
	return b
}

func (b *Builder) WithCustomCollectors(collectors []CollectorFunc) *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, collectors...)
	return b
}

func (b *Builder) WithConfig(cfg Config) *Builder {
	b.aggregator.cfg = cfg
	return b
}

func (b *Builder) WithErrorChan(errorChan chan error) *Builder {
	b.aggregator.errorChan = errorChan
	return b
}

func (b *Builder) Build() (*aggregator, error) {
	if b.err != nil {
		return nil, b.err
	}

	if b.aggregator.errorChan == nil {
		b.aggregator.errorChan = make(chan error)
	}

	return &b.aggregator, b.err
}
