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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemStorage struct {
	mx              sync.Mutex
	l               *zap.SugaredLogger
	Metrics         map[string]models.Metrics
	FileStoragePath string
}

func NewMemStorage(fileStoragePath string, l *zap.SugaredLogger) interfaces.MetricsStorage {
	storage := &MemStorage{
		l:               l,
		Metrics:         make(map[string]models.Metrics),
		FileStoragePath: fileStoragePath,
	}
	return storage
}

func (m *MemStorage) GetAll() map[string]models.Metrics {
	m.mx.Lock()
	defer m.mx.Unlock()
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
	m.mx.Lock()
	defer m.mx.Unlock()
	m.Metrics = metrics
	fmt.Println("Successfully restored values from ", m.FileStoragePath)
	return nil
}
func (m *MemStorage) Save() error {
	m.mx.Lock()
	data, err := json.MarshalIndent(m.Metrics, "", "	")
	m.mx.Unlock()
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
	m.mx.Lock()
	defer m.mx.Unlock()
	if val, ok := m.Metrics[name]; ok {

		return val, nil
	}
	return models.Metrics{}, errors.New("no such metric")
}

func (m *MemStorage) CreateOrUpdate(metric models.Metrics) models.Metrics {

	fmt.Println("Create or update metric")
	var name = metric.Name
	var tp = metric.MType
	m.mx.Lock()
	if _, ok := m.Metrics[name]; ok {
		m.mx.Unlock()
		if tp == "gauge" {
			metric.Delta = nil
			m.mx.Lock()
			m.Metrics[name] = metric
			m.mx.Unlock()
			return metric
		}

		m.mx.Lock()
		defer m.mx.Unlock()
		newDelta := *m.Metrics[name].Delta + *metric.Delta
		m.Metrics[name] = models.Metrics{
			Delta: &newDelta,
			Value: metric.Value,
			MType: tp,
			Name:  name,
		}
		return m.Metrics[name]
	}
	m.mx.Unlock()
	m.mx.Lock()
	m.Metrics[name] = metric
	m.mx.Unlock()
	return metric
}

func (m *MemStorage) Ping() error {
	return errors.New("not implemented")
}

func (m *MemStorage) Remove(name string) error {
	if _, ok := m.Metrics[name]; ok {
		m.mx.Lock()
		delete(m.Metrics, name)
		m.mx.Unlock()
		return nil
	}
	return errors.New("no such metric")

}

type SQLStorage struct {
	l  *zap.SugaredLogger
	db *gorm.DB
}

func NewSQLStorage(URL string, l *zap.SugaredLogger) interfaces.MetricsStorage {

	pgdb, err := sql.Open("pgx", URL)

	if err != nil {
		l.Errorf("Couldn't establish connection with DB on following URL: %s %s", URL, err)
		return NewMemStorage("./save.json", l)
	}

	gormdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: pgdb,
	}), &gorm.Config{})
	if err != nil {
		l.Errorf("Couldn't connect gorm: %s %s", URL, err)
		return NewMemStorage("./save.json", l)
	}
	l.Infof("Successfully connected: %s", URL)

	merr := gormdb.AutoMigrate(&models.Metrics{})
	if merr != nil {
		l.Errorf("Couldn't do migrations:  %s", URL, err)
		return NewMemStorage("./save.json", l)
	}
	l.Infof("Run migrations")
	return &SQLStorage{
		db: gormdb,
		l:  l,
	}
}

func (s *SQLStorage) Ping() error {
	s.l.Infoln("test ping sql")
	sqldb, err := s.db.DB()
	if err != nil {
		s.l.Infof("Couldn't get underlying db from gorm %s", err)
	}
	return sqldb.Ping()
}

func (s *SQLStorage) GetAll() map[string]models.Metrics {
	metrics := make(map[string]models.Metrics)
	s.db.Find(&metrics)
	return metrics
}

func (s *SQLStorage) Restore() error {
	return errors.New("not available / needed in this implementation")
}
func (s *SQLStorage) Save() error {
	return errors.New("not implemented")
}

func (s *SQLStorage) Find(name string) (models.Metrics, error) {
	var metric models.Metrics
	res := s.db.Limit(1).First(&metric, "name = ?", name)
	if res.Error != nil {
		s.l.Errorf("Issue when retrieving metric %v", res.Error)
		return metric, res.Error
	}
	return metric, nil
}

func (s *SQLStorage) CreateOrUpdate(metric models.Metrics) models.Metrics {
	if metric.MType == "counter" {
		existing, err := s.Find(metric.Name)
		if err == nil {
			s.l.Infoln("Creating new")
			newValue := *metric.Delta + *existing.Delta
			metric.Delta = &newValue
			s.db.Model(&models.Metrics{}).Where("name = ?", metric.Name).Update("delta", newValue)
			return metric
		}
	}
	s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(&metric)
	if s.db.Error != nil {
		s.l.Errorf("Something went wrong when creating / upadating value in DB %v", s.db.Error)
	}
	return metric
}
func (s *SQLStorage) Remove(name string) error {
	s.db.Delete(models.Metrics{Name: name})
	if s.db.Error != nil {
		s.l.Errorf("Something went wrong when removing value in DB %v, name: %s", s.db.Error, name)
	}
	return s.db.Error
}
