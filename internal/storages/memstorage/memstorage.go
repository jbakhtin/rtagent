package memstorage

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/rtagent/internal/config"
	"log"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	sync.RWMutex
	items map[string]models.Metric
	snapshot Snapshot
}

func New() MemStorage {
	ctx := context.TODO()
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)

	}

	snapshot, err := NewSnapshot(ctx, cfg)
	if err != nil {
		log.Fatal(err)
		return MemStorage{}
	}

	var metrics map[string]models.Metric
	list, ok := snapshot.Import(ctx)
	if !ok {
		metrics = make(map[string]models.Metric, 0)
	} else {
		metrics = list
	}
	//metrics = make(map[string]models.Metric, 0)

	go snapshot.Exporting(ctx, cfg, &metrics)

	return MemStorage{
		items: metrics,
		snapshot: *snapshot,
	}
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	ms.Lock()
	defer ms.Unlock()

	ms.items[metric.Key()] = metric

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

	result := make(map[string]models.Metric, len(ms.items))

	// Deep copy
	for k, v := range ms.items {
		result[k] = v
	}

	return result, nil
}
