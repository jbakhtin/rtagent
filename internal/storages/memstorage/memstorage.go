package memstorage

import (
	"context"
	"errors"
	"github.com/jbakhtin/rtagent/internal/config"
	"sync"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	snc sync.RWMutex
	snapshot *Snapshot
	items map[string]models.Metric
}

func New(ctx context.Context, cfg config.Config) (MemStorage, error) { // TODO: обработать возврат ошиок
	snapshot, err := NewSnapshot(ctx, cfg)
	if err != nil {
		return MemStorage{}, err
	}

	var metrics map[string]models.Metric
	//list, ok := snapshot.Import(ctx)
	//if !ok {
	//	metrics = make(map[string]models.Metric, 0)
	//} else {
	//	metrics = list
	//}
	metrics = make(map[string]models.Metric, 0)

	//go snapshot.Exporting(ctx, cfg, &metrics)

	return MemStorage{
		items: metrics,
		snapshot: snapshot,
	}, nil
}

func (ms *MemStorage) Set(metric models.Metric) (models.Metric, error) {
	ms.snc.Lock()
	defer ms.snc.Unlock()

	ms.items[metric.Key()] = metric

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metricer, error) {
	ms.snc.Lock()
	defer ms.snc.Unlock()

	if value, ok := ms.items[key]; ok {
		return value, nil
	}

	return nil, errors.New("record not found")
}

func (ms *MemStorage) GetAll() (map[string]models.Metricer, error) {
	ms.snc.Lock()
	defer ms.snc.Unlock()

	result := make(map[string]models.Metricer, len(ms.items))

	// Deep copy
	for k, v := range ms.items {
		result[k] = v
	}

	return result, nil
}
