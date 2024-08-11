package interfaces

import (
	"metal/internal/domain/models"
)

type MetricsStorage interface {
	Find(name string) (models.MetricValue, error)
	CreateOrUpdate(metric models.Metric) models.Metric
	Remove(name string) error
}
