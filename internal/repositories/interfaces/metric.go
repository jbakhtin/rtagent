package interfaces

import "github.com/jbakhtin/rtagent/internal/models"

type MetricRepository interface {
	GetAll() (map[string]models.Metric, error)
	Get(key string) (models.Metric, error)
	Update(models.Metric) (models.Metric, error)
}