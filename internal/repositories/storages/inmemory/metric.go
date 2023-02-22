package inmemory

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
)

type MetricRepository struct {
	memStorage memstorage.MemStorage
}

func NewMetricRepository() (*MetricRepository, error) {
	return &MetricRepository{
		memStorage: memstorage.New(),
	}, nil
}

func (mr *MetricRepository) GetAll() (map[string]models.Metricer, error) {
	var list = make(map[string]models.Metricer, 0)

	gauges, err := mr.memStorage.GetAllGauges()
	if err != nil {
		return nil, err
	}

	counters, err := mr.memStorage.GetAllCounters()
	if err != nil {
		return nil, err
	}

	// TODO: сделать более элегантно
	for k, v := range gauges {
		list[k] = v
	}

	for k, v := range counters {
		list[k] = v
	}

	return list, nil
}

func (mr *MetricRepository) GetCounter(key string) (models.Counter, error) {
	return mr.memStorage.GetCounter(key)
}

func (mr *MetricRepository) GetGauge(key string) (models.Gauge, error) {
	return mr.memStorage.GetGauge(key)
}

func (mr *MetricRepository) UpdateCounter(counter models.Counter) error {
	return mr.memStorage.SetCounter(counter)
}

func (mr *MetricRepository) UpdateGauge(gauge models.Gauge) error {
	return mr.memStorage.SetGauge(gauge)
}
