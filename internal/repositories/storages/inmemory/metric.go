package inmemory

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type Metric struct {
	ctx *context.Context
	memStorage memstorage.MemStorage
}

func NewMetricRepository(ctx *context.Context) *Metric {
	return &Metric{
		ctx: ctx,
		memStorage: memstorage.New(),
	}
}

func (metric *Metric) Get() ([]models.Metric, error) {
	return metric.memStorage.Get()
}

func (metric *Metric) Find(tp, key string) (models.Metric, error) {
	return metric.memStorage.Find(tp, key)
}

func (metric *Metric) Update(tp, k, vl string) (models.Metric, error) {
	return metric.memStorage.Set(tp, k, vl)
}