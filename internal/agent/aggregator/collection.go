package aggregator

import (
	"github.com/jbakhtin/rtagent/internal/types"
	"golang.org/x/exp/maps"
	"sync"
)

type Metrics struct {
	sync.RWMutex
	items map[string]types.Metricer
}

func (m *Metrics) Merge(items map[string]types.Metricer) {
	m.Lock()
	defer m.Unlock()

	maps.Copy(m.items, items)
}

func (m *Metrics) GetAll() map[string]types.Metricer{
	m.Lock()
	defer m.Unlock()

	return m.items
}

func (m *Metrics) Set(key string, metric types.Metricer) {
	m.Lock()
	defer m.Unlock()

	m.items[key] = metric
}
