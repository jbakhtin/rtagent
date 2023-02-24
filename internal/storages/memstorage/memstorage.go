package memstorage

import (
	"errors"
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

	result := make(map[string]models.Metricer, len(ms.items))

	// Deep copy
	for k, v := range ms.items {
		result[k] = v
	}

	return result, nil
}
