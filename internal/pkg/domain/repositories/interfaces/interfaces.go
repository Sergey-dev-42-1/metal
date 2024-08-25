package interfaces

import (
	"metal/internal/pkg/domain/models"
)

type MetricsStorage interface {
	Find(name string) (models.Metric, error)
	GetAll() map[string]models.Metric
	CreateOrUpdate(metric models.Metric) models.Metric
	Remove(name string) error
}
