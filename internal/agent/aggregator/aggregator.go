package aggregator

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"golang.org/x/exp/maps"
	"sync"

	"github.com/jbakhtin/rtagent/internal/types"
)

type Collector func() (map[string]types.Metricer, error)

type Aggregator struct {
	sync.RWMutex
	cfg config.Config
	items  map[string]types.Metricer
	collectors []Collector
	poolCount types.Counter
}

func (a *Aggregator) PoolCount()(map[string]types.Metricer, error) {
	//fmt.Println("PoolCount")
	a.poolCount.Add(1)
	return map[string]types.Metricer{"PollCount": types.Counter(a.poolCount)}, nil
}

func New() (Aggregator, error) {
	aggregator := Aggregator{
		items:  make(map[string]types.Metricer, 0),
		collectors: []Collector{Runtime, Gopsutil, RandomMetric},
	}

	aggregator.collectors = append(aggregator.collectors, aggregator.PoolCount)

	return aggregator, nil
}

func (a *Aggregator) Run(ctx context.Context) (err error) {
	a.Lock()
	defer a.Unlock()

	for _, collector := range a.collectors {
		items, tempErr := collector()
		if tempErr != nil && err == nil {
			err = tempErr
		}

		maps.Copy(a.items, items)

		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregator) AddCollector(collector Collector) error {
	a.collectors = append(a.collectors, collector)
	return nil
}

func (a *Aggregator) GetAll() (map[string]types.Metricer, error) {
	a.Lock()
	defer a.Unlock()

	result := make(map[string]types.Metricer, len(a.items))

	// Deep copy
	for k, v := range a.items {
		result[k] = v
	}

	return result, nil
}
