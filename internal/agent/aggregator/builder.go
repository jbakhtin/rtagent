package aggregator

import "github.com/jbakhtin/rtagent/internal/types"

type Builder struct {
	err        error
	aggregator *aggregator
}

func New() *Builder {
	return &Builder{
		aggregator: &aggregator{
			collection: &Metrics{
				items: make(map[string]types.Metricer, 0),
			},
		},
	}
}

func (b *Builder) WithDefaultCollectors() *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, []CollectorFunc{Runtime, Gopsutil, RandomMetric}...)
	return b
}

func (b *Builder) WithCustomCollector(collector CollectorFunc) *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, collector)
	return b
}

func (b *Builder) WithCustomCollectors(collectors ...CollectorFunc) *Builder {
	b.aggregator.collectors = append(b.aggregator.collectors, collectors...)
	return b
}

func (b *Builder) Build() (*aggregator, error) {
	if b.err != nil {
		return nil, b.err
	}

	return b.aggregator, b.err
}
