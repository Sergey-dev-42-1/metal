package repositories

import (
	"errors"
	"fmt"
	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"
)

type MemStorage struct {
	//дубляжка имени, но зато удобно работать
	Metrics map[string]models.Metric
}

func New() interfaces.MetricsStorage {
	return &MemStorage{
		Metrics: make(map[string]models.Metric),
	}
}

func (m *MemStorage) GetAll() map[string]models.Metric {
	return m.Metrics
}

func (m *MemStorage) Find(name string) (models.Metric, error) {
	if val, ok := m.Metrics[name]; ok {
		return val, nil
	}
	return models.Metric{}, errors.New("no such metric")
}

func (m *MemStorage) CreateOrUpdate(metric models.Metric) models.Metric {

	fmt.Println("Create or update metric")

	var name = metric.Name
	var tp = metric.Type

	if _, ok := m.Metrics[metric.Name]; ok {
		if tp == "gauge" {
			m.Metrics[name] = metric
			return metric
		}
		m.Metrics[name] = models.Metric{
			Value: m.Metrics[name].Value + metric.Value,
			Type:  tp,
			Name:  name,
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
