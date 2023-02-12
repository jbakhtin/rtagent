package interfaces

import "github.com/jbakhtin/rtagent/internal/models"

type MetricRepository interface {
	Get() ([]models.Metric, error)
	Find(t, key string) (models.Metric, error)
	Update(tp, k, vl string) (models.Metric, error)
}