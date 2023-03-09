package inmemory

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/filestorage"
)

type MetricRepository struct {
	memStorage filestorage.FileStorage // TODO: need add interface instead
}

func NewMetricRepository(ctx context.Context, cfg config.Config) (*MetricRepository, error) {
	ms, err := filestorage.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	err = ms.Start(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &MetricRepository{
		memStorage: ms,
	}, nil
}

func (mr *MetricRepository) GetAll() (map[string]models.Metric, error) {
	return mr.memStorage.GetAll()
}

func (mr *MetricRepository) Get(key string) (models.Metric, error) {
	return mr.memStorage.Get(key)
}

func (mr *MetricRepository) Update(metric models.Metric) (models.Metric, error) {
	return mr.memStorage.Set(metric)
}
