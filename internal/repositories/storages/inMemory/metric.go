package inMemory

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/models"
)

type Metric struct {
	ctx *context.Context
}

func (metric *Metric) Get() ([]models.Metric, error) {
	return []models.Metric{}, nil
}

func NewMetricRepository(ctx *context.Context) *Metric {
	return &Metric{
		ctx: ctx,
	}
}

func (metric *Metric) Update(type2, key, value string) (models.Metric, error) {

	return models.Metric{
		Type2: type2 + "_updated",
		Key:   key + "_updated",
		Value: value + "_updated",
	}, nil
}

func (metric *Metric) Create(type2, key, value string) (models.Metric, error) {

	return models.Metric{
		Type2: type2,
		Key:   key,
		Value: value,
	}, nil
}




