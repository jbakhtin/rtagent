package inmemory

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type Metric struct {
	memStorage memstorage.MemStorage
}

func NewMetricRepository() *Metric {
	return &Metric{
		memStorage: memstorage.New(),
	}
}

func (metric *Metric) GetAll() ([]models.Metric, error) {
	return metric.memStorage.GetAll()
}

func (metric *Metric) Get(tp, key string) (models.Metric, error) {
	return metric.memStorage.Get(tp, key)
}

func (metric *Metric) Update(tp, k, vl string) (models.Metric, error) {
	return metric.memStorage.Set(tp, k, vl)
}