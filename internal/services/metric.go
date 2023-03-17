package services

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/storages/filestorage"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Set(models.Metricer) (models.Metricer, error)
}

type MetricService struct {
	repository MetricRepository
}

func NewMetricService(ctx context.Context, cfg config.Config) (*MetricService, error) {
	ms, err := filestorage.New(cfg)
	if err != nil {
		return nil, err
	}

	err = ms.Start(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &MetricService{
		repository: &ms,
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

	switch m := metric.(type) {
	case models.Counter:
		entity, err := ms.repository.Get(m.MKey)
		if err != nil {
			break
		}

		oldMetric, ok := entity.(models.Counter)
		if !ok {
			return nil, err
		}

		m.Add(oldMetric.MValue)
		metric = m
	}

	metric, err = ms.repository.Set(metric)
	if err != nil {
		fmt.Println("Update error: ", err)
		return metric, err
	}

	return metric, nil
}
