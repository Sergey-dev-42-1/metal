package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"
)

type MemStorage struct {
	mx              sync.RWMutex
	Metrics         map[string]models.Metrics
	FileStoragePath string
}

func New(fileStoragePath string) interfaces.MetricsStorage {
	storage := &MemStorage{
		mx:              sync.RWMutex{},
		Metrics:         make(map[string]models.Metrics),
		FileStoragePath: fileStoragePath,
	}
	return storage
}

func (m *MemStorage) GetAll() map[string]models.Metrics {
	return m.Metrics
}

func (m *MemStorage) Restore() error {
	content, err := os.ReadFile(m.FileStoragePath)
	if err != nil {
		fmt.Println("Had an issue when trying to open file", err)
		return err
	}
	// scanner := bufio.NewScanner(file)
	// scanner.Scan()
	// file.Close()
	// fmt.Println(scanner.Text())
	metrics := map[string]models.Metrics{}
	errRead := json.Unmarshal(content, &metrics)
	if errRead != nil {
		fmt.Println("Had an issue when trying to restore saved values", errRead)
		return errRead
	}
	m.Metrics = metrics
	fmt.Println("Successfully restored values from ", m.FileStoragePath)
	return nil
}
func (m *MemStorage) Save() error {
	data, err := json.MarshalIndent(m.Metrics, "", "	")
	if err != nil {
		fmt.Println("Had an issue converting to json", err)
		return err
	}
	writeErr := os.WriteFile(m.FileStoragePath, data, 0666)
	if writeErr != nil {
		fmt.Println("Had an issue when trying to read file", err)
		return err
	}

	fmt.Println("Successfully written to ", m.FileStoragePath)
	return nil
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
			m.mx.Lock()
			m.Metrics[name] = metric
			m.mx.Unlock()
			return metric
		}
		newDelta := *m.Metrics[name].Delta + *metric.Delta

		m.Metrics[name] = models.Metrics{
			Delta: &newDelta,
			Value: metric.Value,
			MType: tp,
			ID:    name,
		}
		return m.Metrics[name]
	}
	m.mx.Lock()
	m.Metrics[name] = metric
	m.mx.Unlock()
	return metric
}

func (m *MemStorage) Remove(name string) error {
	if _, ok := m.Metrics[name]; ok {
		delete(m.Metrics, name)
		return nil
	}
	return errors.New("no such metric")

}
