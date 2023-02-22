package memstorage

import (
	"errors"
	"github.com/jbakhtin/rtagent/internal/models"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	gauges map[string]models.Gauge
	counters map[string]models.Counter
}

func New() MemStorage {
	return MemStorage{
		gauges: make(map[string]models.Gauge, 0),
		counters: make(map[string]models.Counter, 0),
	}
}

func (ms *MemStorage) SetGauge(value models.Gauge) error {
	ms.Lock()
	defer ms.Unlock()

	ms.gauges[value.MKey] = value

	return nil
}

func (ms *MemStorage) SetCounter(value models.Counter) error {
	ms.Lock()
	defer ms.Unlock()

	ms.counters[value.MKey] = value

	return nil
}

func (ms *MemStorage) GetGauge(key string) (models.Gauge, error) {
	ms.Lock()
	defer ms.Unlock()

	if value, ok := ms.gauges[key]; ok {
		return value, nil
	}

	return models.Gauge{}, errors.New("record not found")
}

func (ms *MemStorage) GetCounter(key string) (models.Counter, error) {
	ms.Lock()
	defer ms.Unlock()

	if value, ok := ms.counters[key]; ok {
		return value, nil
	}

	return models.Counter{}, errors.New("record not found")
}

func (ms *MemStorage) GetAllGauges() (map[string]models.Gauge, error) {
	ms.Lock()
	defer ms.Unlock()

	result := make(map[string]models.Gauge, len(ms.gauges))
	result = ms.gauges

	return result, nil
}

func (ms *MemStorage) GetAllCounters() (map[string]models.Counter, error) {
	ms.Lock()
	defer ms.Unlock()

	// deep copy это оно?
	result := make(map[string]models.Counter, len(ms.counters))
	result = ms.counters

	return result, nil
}
