package service

import (
	"fmt"
	"math/rand"
	"metal/internal/agent/application/metrics/interfaces"
	updateService "metal/internal/agent/application/update"
	"metal/internal/pkg/domain/models"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"
)

type MetricsService struct {
	addr           string
	reportInterval int
	pollInterval   int
}

func New(a string, r int, p int) interfaces.MetricsService {
	return &MetricsService{
		addr:           a,
		reportInterval: r,
		pollInterval:   p,
	}
}

var stats []models.Metric
var pollCount int64
var randomValue int64
var collectedMetrics = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle",
	"HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
	"Sys", "TotalAlloc"}

func (s *MetricsService) CollectMemStats() {
	for {
		fmt.Println("Updating metrics")
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		collectedMetrics := createMetricsMap(memStats)
		stats = collectedMetrics
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
	}
}

func (s *MetricsService) SendMemStats() {
	for {
		time.Sleep(time.Duration(s.reportInterval) * time.Second)
		pollCount = 0
		fmt.Println("Saving metrics")
		service := updateService.New(s.addr)
		for _, v := range stats {
			go service.UpdateMetrics(v)
		}
	}
}

func convertToMetrics(m runtime.MemStats) []models.Metric {
	val := reflect.ValueOf(m)
	metrics := []models.Metric{}
	for i := range val.NumField() {
		key := val.Type().Field(i).Name

		keyID := slices.IndexFunc(collectedMetrics, func(metricName string) bool {
			return strings.EqualFold(key, metricName)
		})
		if keyID == -1 {
			continue
		}
		var value models.MetricValue
		fieldValue := val.Field(i)
		fieldType := val.Field(i).Type().Name()

		switch fieldType {
		case "uint64":
			value = models.MetricValue(fieldValue.Uint())
		case "uint32":
			value = models.MetricValue(fieldValue.Uint())
		case "int64":
			value = models.MetricValue(fieldValue.Int())
		case "float64":
			value = models.MetricValue(fieldValue.Float())
		default:
			fmt.Printf("Unsupported value type in struct %s", fieldType)
			continue
		}

		metrics = append(metrics, models.Metric{
			Value: value,
			Type:  "gauge",
			Name:  key,
		})

	}
	return metrics
}

func createMetricsMap(m runtime.MemStats) []models.Metric {

	pollCount++
	randomValue = rand.Int63n(1000000)

	metricsMap := []models.Metric{}

	metricsMap = append(metricsMap, models.Metric{
		Value: models.MetricValue(pollCount),
		Type:  "counter",
		Name:  "PollCount",
	}, models.Metric{
		Value: models.MetricValue(randomValue),
		Type:  "gauge",
		Name:  "RandomValue",
	})
	metricsMap = append(metricsMap, convertToMetrics(m)...)
	return metricsMap
}
