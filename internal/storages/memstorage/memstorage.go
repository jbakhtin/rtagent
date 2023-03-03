package memstorage

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/jbakhtin/rtagent/internal/config"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	mx       *sync.RWMutex
	snapshot Snapshot
	items    map[string]models.Metric
}

func New(ctx context.Context, cfg config.Config) (MemStorage, error) {
	snapshot, err := NewSnapshot(ctx, cfg)
	if err != nil {
		log.Fatal(err)
		return MemStorage{}, err
	}

	var metrics map[string]models.Metric
	list, ok := snapshot.Import(ctx)
	if !ok {
		metrics = make(map[string]models.Metric, 0)
	} else {
		metrics = list
	}

	go snapshot.Exporting(ctx, cfg, &metrics)

	return MemStorage{
		items:    metrics,
		snapshot: *snapshot,
		mx:       &sync.RWMutex{},
	}, nil
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	ms.items[metric.Key()] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metric, error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	if value, ok := ms.items[key]; ok {
		return value, nil
	}

	return models.Metric{}, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metric, error) {
	ms.mx.Lock()
	defer ms.mx.Unlock()

	result := make(map[string]models.Metric, len(ms.items))

	// Deep copy
	for k, v := range ms.items {
		result[k] = v
	}

	return result, nil
}
