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

func (ms *MetricService) Get() ([]models.Metric, error){
	metrics, err := ms.repository.Get()

	if err != nil {
		fmt.Println("Update error: ", err)
	}

	return metrics, nil
}

func (ms *MetricService) Update(type2, key, value string) (models.Metric, error) {
	metric, err := ms.repository.Update(type2, key, value)

	if err != nil {
		fmt.Println("Update error: ", err)
	}

	return metric, nil
}

func (ms *MetricService) Create(type2, key, value string) (models.Metric, error) {
	metric, err := ms.repository.Create(type2, key, value)

	if err != nil {
		fmt.Println("Crated error: ", err)
	}

	return metric, nil
}
