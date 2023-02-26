package services

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/jbakhtin/rtagent/internal/types"

	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
)

type MetricService struct {
	repository interfaces.MetricRepository
}

func NewMetricService() (*MetricService, error) {
	// Инициализируе нужное хранилище
	repository, err := inmemory.NewMetricRepository()
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

func (ms *MetricService) GetAll() (map[string]models.Metricer, error) {
	metrics, err := ms.repository.GetAll()
	if err != nil {
		fmt.Println("Get error: ", err)
		return metrics, err
	}

	return metrics, nil
}

func (ms *MetricService) Update(metric models.Metricer) (models.Metricer, error) {
	var err error

	m, _ := metric.(models.Metric)
	switch Value := m.MValue.(type) { // TODO: подумать как сделать код элегантнее
	case types.Counter:
		entity, err := ms.repository.Get(m.MKey)
		if err != nil {
			break
		}

		oldMetric, ok := entity.(models.Metric)
		if !ok {
			return nil, err
		}

		counter := oldMetric.MValue.(types.Counter)

		Value.Add(counter)
		m.MValue = Value
		metric = m
	}

	metric, err = ms.repository.Update(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
