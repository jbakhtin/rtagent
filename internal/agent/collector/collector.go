package collector

import (
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
	"sync"
)

type Collector struct {
	Mx     *sync.RWMutex
	Items map[string]types.Metricer
	Logger *zap.Logger
}

func NewCollector() (Collector, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return Collector{}, err
	}

	return Collector{
		Items:  make(map[string]types.Metricer, 0),
		Mx:     &sync.RWMutex{},
		Logger: logger,
	}, nil
}

func (ms *Collector) Set(key string, metric types.Metricer) (types.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	ms.Items[key] = metric

	return metric, nil
}

func (ms *Collector) GetAll() (map[string]types.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	result := make(map[string]types.Metricer, len(ms.Items))

	// Deep copy
	for k, v := range ms.Items {
		result[k] = v
	}

	return result, nil
}


