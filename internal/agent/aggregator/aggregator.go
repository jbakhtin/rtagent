package aggregator

import (
	"context"
	"golang.org/x/exp/maps"
	"sync"
	"time"

	"github.com/jbakhtin/rtagent/internal/types"
)

type CollectorFunc func() (map[string]types.Metricer, error)

type aggregator struct {
	sync.RWMutex
	cfg Config
	items  map[string]types.Metricer
	collectors []CollectorFunc
	poolCount types.Counter
	errorChan chan error
	doneChan chan struct{}
}

func (a *aggregator) poolCountCollector()(map[string]types.Metricer, error) {
	a.poolCount.Add(1)
	return map[string]types.Metricer{"PollCount": types.Counter(a.poolCount)}, nil
}

func (a *aggregator) run(ctx context.Context) (err error) {
	a.Lock()
	defer a.Unlock()

	for _, collector := range a.collectors {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			items, tempErr := collector()
			if tempErr != nil && err == nil {
				err = tempErr
			}

			maps.Copy(a.items, items)

			if err != nil {
				return err
			}
		}
	}

	return
}

func (a *aggregator) Run(ctx context.Context) error {
	ticker := time.NewTicker(a.cfg.GetPollInterval())
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go func() {
				err := a.run(ctx)
				if err != nil {
					a.errorChan <- err
				}
			}()
		}
	}
}

func (a *aggregator) GetAll() (map[string]types.Metricer, error) {
	a.Lock()
	defer a.Unlock()

	result := make(map[string]types.Metricer, len(a.items))

	// Deep copy
	for k, v := range a.items {
		result[k] = v
	}

	return result, nil
}

func (a *aggregator) Err() chan error {
	return a.errorChan
}
