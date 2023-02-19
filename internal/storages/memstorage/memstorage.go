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

func (memStorage *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	memStorage.items[metric.MKey] = metric

	return metric, nil
}

func (memStorage *MemStorage) Get(key string) (models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	if value, ok := memStorage.items[key]; ok {
		return value, nil
	}

	return models.Metric{}, errors.New("record not found")
}

func (memStorage *MemStorage) GetAll() (map[string]models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	return memStorage.items, nil
}
