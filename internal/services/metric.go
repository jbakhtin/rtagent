package services

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/infile"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"

	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
)

type MetricService struct {
	repository interfaces.MetricRepository
}

func NewMetricService(ctx context.Context, cfg config.Config) (*MetricService, error) {
	repository, err := infile.NewMetricRepository(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &MetricService{
		repository: repository,
	}, nil
}

func (ms *MetricService) Get(key string) (models.Metricer, error) {
	metric, err := ms.repository.Get(key)
	if err != nil {
		fmt.Println("Find error: ", err)
		return metric, err
	}

	return metric, nil
}

func (ms *MetricService) GetAll() (map[string]models.Metric, error) {
	metrics, err := ms.repository.GetAll()
	if err != nil {
		fmt.Println("Get error: ", err)
		return metrics, err
	}

	return metrics, nil
}

func (ms *MetricService) Update(metric models.Metric) (models.Metric, error) {
	var err error

	switch metric.MType {
	case types.CounterType:
		entity, err := ms.repository.Get(metric.MKey)
		if err != nil {
			break
		}

		metric.Delta.Add(*entity.Delta)
	}

	metric, err = ms.repository.Update(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
