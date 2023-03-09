package memstorage

import (
	"context"
	"errors"
	"sync"

	"github.com/jbakhtin/rtagent/internal/config"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	Mx       *sync.RWMutex
	Items    map[string]models.Metric
}

func NewMemStorage(ctx context.Context, cfg config.Config) (MemStorage, error) {
	return MemStorage{
		Items:    make(map[string]models.Metric, 0),
		Mx:       &sync.RWMutex{},
	}, nil
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	ms.Items[metric.Key()] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metric, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	if value, ok := ms.Items[key]; ok {
		return value, nil
	}

	return models.Metric{}, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metric, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	result := make(map[string]models.Metric, len(ms.Items))

	// Deep copy
	for k, v := range ms.Items {
		result[k] = v
	}

	return result, nil
}
