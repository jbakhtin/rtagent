package memstorage

import (
	"errors"
	"github.com/jbakhtin/rtagent/internal/models"
	"reflect"
	"sync"
)

type entity interface {
	Type() string
}
type gauge float64
type counter int64

func (g gauge) Type() string {
	return reflect.TypeOf(g).Name()
}

func (c counter) Type() string {
	return reflect.TypeOf(c).Name()
}

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

func (memStorage *MemStorage) Get(k string) (models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	if value, ok := memStorage.items[k]; ok {
		return value, nil
	}

	return models.Metric{}, errors.New("record not found")
}

func (memStorage *MemStorage) GetAll() (map[string]models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	return memStorage.items, nil
}
