package interfaces

import (
	"github.com/jbakhtin/rtagent/internal/models"
)

type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	UpdateCounter(counter models.Counter) error
	UpdateGauge(gauge models.Gauge) error
	GetCounter(key string) (models.Counter, error)
	GetGauge(key string) (models.Gauge, error)
}
