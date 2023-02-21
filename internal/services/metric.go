package services

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"

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

func (ms *MetricService) Get(key string) (models.Metric, error) {
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
	case models.GaugeType: // do nothing
	case models.CounterType:
		oldMetric, _ := ms.repository.Get(metric.MKey)
		metric.MCounter += oldMetric.MCounter
	}

	metric, err = ms.repository.Update(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
