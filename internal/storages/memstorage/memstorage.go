package memstorage

import (
	"context"
	"errors"
	"github.com/jbakhtin/rtagent/internal/config"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	sync.RWMutex
	snapshot *Snapshot
	items map[string]models.Metric
}

func New(ctx context.Context, cfg config.Config) MemStorage { // TODO: обработать возврат ошиок
	snapshot, err := NewSnapshot(ctx, cfg)
	if err != nil {
		return MemStorage{}
	}

	metrics, err := snapshot.Import(ctx)
	if err != nil {
		return MemStorage{}
	}

	if metrics == nil {
		metrics = make(map[string]models.Metric, 0)
	}

	go snapshot.Exporting(ctx, cfg, &metrics)

	return MemStorage{
		items: metrics,
		snapshot: snapshot,
	}
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
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
