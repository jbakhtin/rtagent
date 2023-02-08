package interfaces

import (
	"github.com/jbakhtin/rtagent/internal/models"
)

type MetricRepository interface {
	Get() ([]models.Metric, error)
	Create(type2, key, value string) (models.Metric, error)
	Update(type2, key, value string) (models.Metric, error)
}