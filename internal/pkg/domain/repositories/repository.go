package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"


	"go.uber.org/zap"
)

type MemStorage struct {
	mx              sync.RWMutex
	l               *zap.SugaredLogger
	Metrics         map[string]models.Metrics
	FileStoragePath string
}

func NewMemStorage(fileStoragePath string, l *zap.SugaredLogger) interfaces.MetricsStorage {
	storage := &MemStorage{
		mx:              sync.RWMutex{},
		l:               l,
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
		m.mx.Lock()
		if tp == "gauge" {
			metric.Delta = nil
			
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
		m.mx.Unlock()
		return m.Metrics[name]
	}
	m.mx.Lock()
	m.Metrics[name] = metric
	m.mx.Unlock()
	return metric
}

func (m *MemStorage) Ping() error {
	return errors.New("Not implemented")
}

func (m *MemStorage) Remove(name string) error {
	if _, ok := m.Metrics[name]; ok {
		delete(m.Metrics, name)
		return nil
	}
	return errors.New("no such metric")

}

type SQLStorage struct {
	l  *zap.SugaredLogger
	db *sql.DB
}

func NewSQLStorage(URL string, l *zap.SugaredLogger) interfaces.MetricsStorage {
	db, err := sql.Open("pgx", URL)
	if err != nil {
		l.Errorf("Couldn't establish connection with DB on following URL: %s %s", URL, err)
		return NewMemStorage("./save.json", l)
	}

	return &SQLStorage{
		db: db,
		l:  l,
	}
}

func (s *SQLStorage) Ping() error {
	s.l.Infoln("test ping sql")
	err := s.db.Ping()
	if (err != nil) {
		s.l.Infof("test ping sql %s", err)
	}
	return s.db.Ping()
}

func (s *SQLStorage) GetAll() map[string]models.Metrics {
	return make(map[string]models.Metrics)
}

func (s *SQLStorage) Restore() error {
	return errors.New("not available / needed in this implementation")
}
func (s *SQLStorage) Save() error {
	return errors.New("not implemented")
}

func (s *SQLStorage) Find(name string) (models.Metrics, error) {
	var ptr models.Metrics
	return ptr, errors.New("not implemented")
}

func (s *SQLStorage) CreateOrUpdate(metric models.Metrics) models.Metrics {
	var ptr models.Metrics
	return ptr
}

func (s *SQLStorage) Remove(name string) error {
	return errors.New("not implemented")
}
