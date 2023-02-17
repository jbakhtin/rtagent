package interfaces

import "github.com/jbakhtin/rtagent/internal/models"

type MetricRepository interface {
	GetAll() ([]models.Metric, error)
	Get(t, key string) (models.Metric, error)
	Update(tp, k, vl string) (models.Metric, error)
}