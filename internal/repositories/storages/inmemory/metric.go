package inmemory

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type MetricRepository struct {
	memStorage memstorage.MemStorage
}

func NewMetricRepository() (*MetricRepository, error) {
	return &MetricRepository{
		memStorage: memstorage.New(),
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
