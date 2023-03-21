package memstorage

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/jbakhtin/rtagent/internal/config"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MemStorage struct {
	Mx     *sync.RWMutex
	Items  map[string]models.Metricer
	Logger *zap.Logger
}

func NewMemStorage(cfg config.Config) (MemStorage, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return MemStorage{}, err
	}

	return MemStorage{
		Items:  make(map[string]models.Metricer, 0),
		Mx:     &sync.RWMutex{},
		Logger: logger,
	}, nil
}

func (ms *MemStorage) Set(metric models.Metricer) (models.Metricer, error) {
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

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

func (ms *MemStorage) SetBatch(ctx context.Context, metrics []models.Metricer) ([]models.Metricer, error){
	ms.Mx.Lock()
	defer ms.Mx.Unlock()

	for _, v := range metrics {
		ms.Items[v.Key()] = v
	}

	return metrics, nil
}
