package inmemory

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type MetricRepository struct {
	memStorage memstorage.MemStorage
}

func NewMetricRepository(ctx context.Context, cfg config.Config) (*MetricRepository, error) {
	ms, err := memstorage.NewMemStorage(cfg)
	if err != nil {
		return nil, err
	}

	return &MetricRepository{
		memStorage: ms,
	}, nil
}

func (mr *MetricRepository) GetAll() (map[string]models.Metricer, error) {
	return mr.memStorage.GetAll()
}

func (mr *MetricRepository) Get(key string) (models.Metricer, error) {
	return mr.memStorage.Get(key)
}

func (mr *MetricRepository) Update(metric models.Metricer) (models.Metricer, error) {
	return mr.memStorage.Set(metric)
}
