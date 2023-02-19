package inmemory

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type MetricRepository struct {
	memStorage memstorage.MemStorage
}

func NewMetricRepository() *MetricRepository {
	return &MetricRepository{
		memStorage: memstorage.New(),
	}
}

func (metricRepo *MetricRepository) GetAll() (map[string]models.Metric, error) {
	return metricRepo.memStorage.GetAll()
}

func (metricRepo *MetricRepository) Get(key string) (models.Metric, error) {
	return metricRepo.memStorage.Get(key)
}

func (metricRepo *MetricRepository) Update(metric models.Metric) (models.Metric, error) {
	return metricRepo.memStorage.Set(metric)
}
