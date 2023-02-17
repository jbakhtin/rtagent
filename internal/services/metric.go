package services

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"strconv"
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

func (ms *MetricService) Get(tp, key string) (models.Metric, error){
	metric, err := ms.repository.Get(tp, key)

	if err != nil {
		fmt.Println("Find error: ", err)
	}

	return metric, nil
}

func (ms *MetricService) GetAll() ([]models.Metric, error){
	metrics, err := ms.repository.GetAll()

	if err != nil {
		fmt.Println("Get error: ", err)
	}

	return metrics, nil
}

func (ms *MetricService) Update(tp, k, vl string) (models.Metric, error) {
	var metric models.Metric
	var err error

	if tp == "gauge" {
		metric, err = ms.repository.Update(tp, k, vl)
	} else if tp == "counter" {
		metric, err = ms.repository.Get(tp, k)
		if err != nil {
			metric, err = ms.repository.Update(tp, k, vl)
		} else {
			int1, _ := strconv.Atoi(metric.Value())
			int2, _ := strconv.Atoi(vl)

			int3 := int1 + int2

			metric, err = ms.repository.Update(tp, k, strconv.Itoa(int3))
		}
	}

	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
