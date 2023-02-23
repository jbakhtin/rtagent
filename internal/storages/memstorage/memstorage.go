package memstorage

import (
	"errors"
	"golang.org/x/exp/maps"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	sync.RWMutex
	items map[string]models.Metricer
}

func New() MemStorage {
	return MemStorage{
		items: make(map[string]models.Metricer, 0),
	}
}

func (ms *MemStorage) Set(metric models.Metricer) (models.Metricer, error) {
	ms.Lock()
	defer ms.Unlock()

	ms.items[metric.Key()] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metricer, error) {
	ms.Lock()
	defer ms.Unlock()

	if value, ok := ms.items[key]; ok {
		return value, nil
	}

	return nil, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metricer, error) {
	ms.Lock()
	defer ms.Unlock()

	// Deep copy
	metrics := make(map[string]models.Metricer, len(ms.items))
	maps.Copy(metrics, ms.items)

	return metrics, nil
}
