package aggregator

import (
	"github.com/jbakhtin/rtagent/internal/types"
	"golang.org/x/exp/maps"
	"sync"
)

type Metrics struct {
	items map[string]types.Metricer
	sync.RWMutex
}

func (m *Metrics) Merge(items map[string]types.Metricer) {
	m.Lock()
	defer m.Unlock()

	maps.Copy(m.items, items)
}

func (m *Metrics) GetAll() map[string]types.Metricer {
	m.Lock()
	defer m.Unlock()

	// Deep copy
	result := make(map[string]types.Metricer, len(m.items))
	for k, v := range m.items {
		result[k] = v
	}

	return result
}

func (m *Metrics) Set(key string, metric types.Metricer) {
	m.Lock()
	defer m.Unlock()

	m.items[key] = metric
}
