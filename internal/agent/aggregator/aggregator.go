package aggregator

import (
	"context"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/internal/types"
)

type CollectorFunc func() (map[string]types.Metricer, error)

type aggregator struct {
	sync.RWMutex
	cfg Config
	collection  Metrics
	collectors []CollectorFunc
	poolCount types.Counter
	errorChan chan error
	doneChan chan struct{}
}

func (a *aggregator) poolCountCollector()(map[string]types.Metricer, error) {
	a.poolCount.Add(1)
	return map[string]types.Metricer{"PollCount": types.Counter(a.poolCount)}, nil
}

func (a *aggregator) run(ctx context.Context) {
	a.Lock()
	defer a.Unlock()

	for i, _ := range a.collectors {
		select {
		case <-ctx.Done():
			return
		default:
			go func(index int) {
				metrics, err := a.collectors[index]()
				if err != nil {
					a.errorChan<- err
					return
				}

				for key, metric := range metrics {
					a.collection.Set(key, metric)
				}
			}(i)
		}
	}

	return
}

func (a *aggregator) Run(ctx context.Context) {
	ticker := time.NewTicker(a.cfg.GetPollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.run(ctx)
		}
	}
}

func (a *aggregator) GetAll() (map[string]types.Metricer, error) {
	a.Lock()
	defer a.Unlock()

	result := make(map[string]types.Metricer, len(a.collection.GetAll()))

	// Deep copy
	for k, v := range a.collection.GetAll() {
		result[k] = v
	}

	return result, nil
}

func (a *aggregator) Err() chan error {
	return a.errorChan
}
