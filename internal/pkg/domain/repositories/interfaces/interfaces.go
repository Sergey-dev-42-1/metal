package interfaces

import (
	"metal/internal/pkg/domain/models"
)

type MetricsStorage interface {
	Find(name string) (models.Metrics, error)
	GetAll() map[string]models.Metrics
	CreateOrUpdate(metric models.Metrics) models.Metrics
	Remove(name string) error
}
