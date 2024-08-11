package service

import (
	"metal/internal/domain/models"
	"metal/internal/domain/repositories/interfaces"
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

func CreateOrUpdateMetric(metric models.Metric) models.Metric {
	return db.CreateOrUpdate(metric)
}

func FindMetric(name string) (models.MetricValue, error) {
	return db.Find(name)
}


func SetStorage(ms interfaces.MetricsStorage) {
	db = ms
}
