package interfaces

import "github.com/jbakhtin/rtagent/internal/models"

type MetricRepository interface {
	GetAll() (map[string]models.Metricer, error)
	Get(key string) (models.Metricer, error)
	Update(models.Metric) (models.Metric, error)
}
