package aggregator

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/types"
	"golang.org/x/sync/errgroup"
	"sync"
)

type CollectorFunc func() (map[string]types.Metricer, error)

type aggregator struct {
	collectors []CollectorFunc
	collection Metrics
	sync.RWMutex
	poolCount  types.Counter
}

func (a *aggregator) poolCountCollector() (map[string]types.Metricer, error) {
	a.poolCount.Add(1)
	return map[string]types.Metricer{"PollCount": types.Counter(a.poolCount)}, nil
}

func (a *aggregator) Pool(ctx context.Context) error {
	a.Lock()
	defer a.Unlock()

	eg := errgroup.Group{}

	for i := range a.collectors {
		i := i
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			eg.Go(func() error {
				metrics, err := a.collectors[i]()
				if err != nil {
					return err
				}

				for key, metric := range metrics {
					a.collection.Set(key, metric)
				}

				return nil
			})
		}
	}

	return eg.Wait()
}

func (a *aggregator) GetAll() map[string]types.Metricer {
	a.Lock()
	defer a.Unlock()

	result := a.collection.GetAll()

	return result
}
