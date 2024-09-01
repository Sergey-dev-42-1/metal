package service

import (
	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"
)

var db interfaces.MetricsStorage

// type MetricsService struct {
// 	repo interfaces.MetricsStorage
// }

// func New(r interfaces.MetricsStorage) *MetricsService {
// 	return &MetricsService{
// 		repo: r,
// 	}
// }

func CreateOrUpdateMetric(metric models.Metrics) models.Metrics {
	return db.CreateOrUpdate(metric)
}

func GetAllMetrics() map[string]models.Metrics {
	return db.GetAll()
}

func FindMetric(name string) (models.Metrics, error) {
	return db.Find(name)
}

func SetStorage(ms interfaces.MetricsStorage) {
	db = ms
}
