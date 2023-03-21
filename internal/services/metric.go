package services

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/storages/filestorage"
	"github.com/jbakhtin/rtagent/internal/storages/postgres"

	"github.com/jbakhtin/rtagent/internal/models"
)

type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Set(models.Metricer) (models.Metricer, error)
	SetBatch(context.Context, []models.Metricer) ([]models.Metricer, error)
}

type MetricService struct {
	repository MetricRepository
}

func NewMetricService(ctx context.Context, cfg config.Config) (*MetricService, error) {
	if cfg.DatabaseDSN != "" {
		ms, err := postgres.New("pgx", cfg) //TODO: move to cfg
		if err != nil {
			return nil, err
		}

		return &MetricService{
			repository: &ms,
		}, nil
	}

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
		return metric, err
	}

	return metric, nil
}

func (ms *MetricService) GetAll() (map[string]models.Metricer, error) {
	metrics, err := ms.repository.GetAll()
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (ms *MetricService) Update(metric models.Metricer) (models.Metricer, error) {
	var err error

	metric, err = ms.repository.Set(metric)
	if err != nil {
		return metric, err
	}

	return metric, nil
}

func (ms *MetricService) UpdateBatch(metrics []models.Metricer) ([]models.Metricer, error) {
	var err error

	metrics, err = ms.repository.SetBatch(context.TODO(), metrics)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}
