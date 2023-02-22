package services

import (
	"fmt"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"

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

func (ms *MetricService) GetCounter(key string) (models.Counter, error) {
	metric, err := ms.repository.GetCounter(key)
	if err != nil {
		fmt.Println("Find error: ", err)
		return metric, err
	}

	return metric, nil
}

func (ms *MetricService) GetGauge(key string) (models.Gauge, error) {
	metric, err := ms.repository.GetGauge(key)
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

func (ms *MetricService) UpdateCounter(metric models.Counter) (models.Counter, error) {
	oldMetric, _ := ms.repository.GetCounter(metric.MKey)
	metric.MValue += oldMetric.MValue

	err := ms.repository.UpdateCounter(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}

func (ms *MetricService) UpdateGauge(metric models.Gauge) (models.Gauge, error) {
	err := ms.repository.UpdateGauge(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
