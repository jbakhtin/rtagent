package services

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
)

type MetricService struct {
	repository interfaces.MetricRepository
	ctx *context.Context
}

func NewMetricService (ctx *context.Context, repository interfaces.MetricRepository) *MetricService {
	return &MetricService{
		repository: repository,
		ctx: ctx,
	}
}

func (ms *MetricService) Find(tp, key string) (models.Metric, error){
	metric, err := ms.repository.Find(tp, key)

	if err != nil {
		fmt.Println("Find error: ", err)
	}

	return metric, nil
}

func (ms *MetricService) Get() ([]models.Metric, error){
	metrics, err := ms.repository.Get()

	if err != nil {
		fmt.Println("Update error: ", err)
	}

	return metrics, nil
}

func (ms *MetricService) Update(tp, k, vl string) (models.Metric, error) {
	metric, err := ms.repository.Update(tp, k, vl)

	if err != nil {
		fmt.Println("Update error: ", err)
	}

	return metric, nil
}
