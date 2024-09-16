package service

import (
	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"

	"go.uber.org/zap"
)

type MetricsService struct {
	repo interfaces.MetricsStorage
	l    *zap.SugaredLogger
}

func New(r interfaces.MetricsStorage, l *zap.SugaredLogger) *MetricsService {
	return &MetricsService{
		repo: r,
		l:    l,
	}
}

func (s *MetricsService) CreateOrUpdateMetric(metric models.Metrics) models.Metrics {
	return s.repo.CreateOrUpdate(metric)
}

func (s *MetricsService) CreateOrUpdateMetricBatch(metrics []models.Metrics) error {
	return s.repo.CreateOrUpdateBatch(metrics)
}

func (s *MetricsService) GetAllMetrics() map[string]models.Metrics {
	return s.repo.GetAll()
}

func (s *MetricsService) FindMetric(name string) (models.Metrics, error) {
	return s.repo.Find(name)
}
func (s *MetricsService) Ping() error {
	return s.repo.Ping()
}
