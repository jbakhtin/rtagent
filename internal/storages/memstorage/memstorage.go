package memstorage

import (
	"errors"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	sync.RWMutex
	items map[string]models.Metric
}

func New() MemStorage {
	return MemStorage{
		items: make(map[string]models.Metric, 0),
	}
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	ms.Lock()
	defer ms.Unlock()

	ms.items[metric.MKey] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metric, error) {
	ms.Lock()
	defer ms.Unlock()

	if value, ok := ms.items[key]; ok {
		return value, nil
	}

	return models.Metric{}, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metric, error) {
	ms.Lock()
	defer ms.Unlock()

	return ms.items, nil
}
