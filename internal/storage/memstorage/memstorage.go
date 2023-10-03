package memstorage

import (
	"errors"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	Mx     *sync.RWMutex
	Items  map[string]models.Metricer
}

func (ms *MemStorage) Set(metric models.Metricer) (models.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	switch m := metric.(type) {
	case models.Counter:
		entity := ms.Items[metric.Key()]

		oldMetric, ok := entity.(models.Counter)
		if !ok {
			break
		}

		m.Add(oldMetric.MValue)
		metric = m
	}

	ms.Items[metric.Key()] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	if value, ok := ms.Items[key]; ok {
		return value, nil
	}

	return nil, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	result := make(map[string]models.Metricer, len(ms.Items))

	// Deep copy
	for k, v := range ms.Items {
		result[k] = v
	}

	return result, nil
}

func (ms *MemStorage) SetBatch(metrics []models.Metricer) ([]models.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	for _, v := range metrics {
		ms.Items[v.Key()] = v
	}

	return metrics, nil
}

func (ms *MemStorage) TestPing() error {
	return nil
}
