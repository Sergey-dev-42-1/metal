package repositories

import (
	"errors"
	"fmt"
	"metal/internal/pkg/domain/models"
)

type MemStorage struct {
	Metrics map[string]models.MetricValue
}

func New() *MemStorage {
	return &MemStorage{
		Metrics: make(map[string]models.MetricValue),
	}
}

func (m *MemStorage) Find(name string) (models.MetricValue, error) {
	if val, ok := m.Metrics[name]; ok {
		return val, nil
	}
	return 0, errors.New("no such metric")
}

func (m *MemStorage) CreateOrUpdate(metric models.Metric) models.Metric {
	
	fmt.Println("Create or update metric")
	name, value := metric.Name, metric.Value
	if metric.Type == "counter" {
		value = models.MetricValue(value.ToInt64())
	}
	if _, ok := m.Metrics[metric.Name]; ok {
		if metric.Type == "gauge" {
			m.Metrics[name] = value
			return models.Metric{Name: name, Value: value}
		}
		m.Metrics[name] += value
		return models.Metric{Name: name, Value: m.Metrics[name]}
	}
	m.Metrics[name] = value
	return models.Metric{Name: name, Value: value}
}

func (m *MemStorage) Remove(name string) error {
	if _, ok := m.Metrics[name]; ok {
		delete(m.Metrics, name)
		return nil
	}
	return errors.New("no such metric")

}
