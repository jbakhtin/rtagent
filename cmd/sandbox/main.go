package main

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
)

func main() {
	metrics := make([]models.Metricer, 0)

	gauge := models.Gauge{
		Description: models.Description{
			MType: "gauge",
			MKey: "key 1",
		},
		MValue: types.Gauge(12.3),
	}

	counter := models.Counter{
		Description: models.Description{
			MType: "counter",
			MKey: "key 2",
		},
		MValue: types.Counter(12),
	}

	metrics = append(metrics, gauge)
	metrics = append(metrics, counter)

	for _, v := range metrics {
		switch v.(type) {
		case models.Gauge:
			metric, _ := v.(models.Gauge)
			metric.Type()

		case models.Counter:
			metric, _ := v.(models.Counter)
			metric.Increment()
		}
	}
}
