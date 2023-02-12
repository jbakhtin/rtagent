package memstorage

import (
	"errors"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"golang.org/x/exp/slices"
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	items []models.Metric
}

func New() MemStorage {
	return MemStorage{
		items: make([]models.Metric, 0),
	}
}

func (memStorage *MemStorage) Set(tp, k, vl string) (models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	metric := models.NewMetric(tp, k, vl)

	idx := slices.IndexFunc(memStorage.items, func(c models.Metric) bool { return c.Type() == "k" })

	if idx == -1 {
		memStorage.items = append(memStorage.items, metric)
	} else {
		memStorage.items[idx] = metric
	}

	return metric, nil
}

func (memStorage *MemStorage) Find(tp, k string) (models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	idx := slices.IndexFunc(memStorage.items, func(c models.Metric) bool {
		return c.Type() == tp && c.Key() == k
	})

	if idx == -1 {
		return models.Metric{}, errors.New("metric not found")
	}

	return memStorage.items[idx], nil
}

func (memStorage *MemStorage) Get() ([]models.Metric, error) {
	memStorage.Lock()
	defer memStorage.Unlock()

	fmt.Println(memStorage.items)

	return memStorage.items, nil
}
