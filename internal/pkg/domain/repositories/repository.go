package repositories

import (
	"errors"
	"fmt"

	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"
)

type MemStorage struct {
	Metrics map[string]models.Metrics
}

func New() interfaces.MetricsStorage {
	return &MemStorage{
		Metrics: make(map[string]models.Metrics),
	}
}

func (m *MemStorage) GetAll() map[string]models.Metrics {
	return m.Metrics
}

func (m *MemStorage) Find(name string) (models.Metrics, error) {
	if val, ok := m.Metrics[name]; ok {
		return val, nil
	}
	return models.Metrics{}, errors.New("no such metric")
}

func (m *MemStorage) CreateOrUpdate(metric models.Metrics) models.Metrics {
	fmt.Println("Create or update metric")
	var name = metric.ID
	var tp = metric.MType

	if _, ok := m.Metrics[name]; ok {
		if tp == "gauge" {
			metric.Delta = nil
			m.Metrics[name] = metric
			return metric
		}
		fmt.Println("test",*m.Metrics[name].Delta, *metric.Delta)
		newDelta := *m.Metrics[name].Delta + *metric.Delta
		m.Metrics[name] = models.Metrics{
			Delta: &newDelta,
			Value: nil,
			MType: tp,
			ID:    name,
		}
		return m.Metrics[name]
	}
	m.Metrics[name] = metric
	return metric
}

func (m *MemStorage) Remove(name string) error {
	if _, ok := m.Metrics[name]; ok {
		delete(m.Metrics, name)
		return nil
	}
	return errors.New("no such metric")

}
