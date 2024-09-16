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
	batching       bool
}

func New(a string, r int, p int, b bool) interfaces.MetricsService {
	return &MetricsService{
		addr:           a,
		reportInterval: r,
		pollInterval:   p,
		batching:       b,
	}
}

var stats []models.Metrics
var pollCount int64
var randomValue float64
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
		if s.batching {
			//but why tho?
			stats = append(stats, collectedMetrics...)
		} else {
			stats = collectedMetrics
		}
		time.Sleep(time.Duration(s.pollInterval) * time.Second)
	}
}

func (s *MetricsService) SendMemStats() {
	for {
		time.Sleep(time.Duration(s.reportInterval) * time.Second)
		service := updateService.New(s.addr)
		if s.batching {
			// don't think there is need to run goroutine
			service.UpdateMetricsJSONBatch(stats)
			stats = make([]models.Metrics, 0)
			pollCount = 0
			continue
		}
		for _, v := range stats {
			if v.Name == "PollCount" {
				go func() {
					service.UpdateMetricsJSON(v)
					pollCount = 0
				}()
			}
			go service.UpdateMetricsJSON(v)
		}
	}
}

func convertToMetrics(m runtime.MemStats) []models.Metrics {
	val := reflect.ValueOf(m)
	metrics := []models.Metrics{}
	for i := range val.NumField() {
		key := val.Type().Field(i).Name

		keyID := slices.IndexFunc(collectedMetrics, func(metricName string) bool {
			return strings.EqualFold(key, metricName)
		})
		if keyID == -1 {
			continue
		}
		metric := models.Metrics{}
		fieldValue := val.Field(i)
		fieldType := val.Field(i).Type().Name()

		switch fieldType {
		case "uint64":
			val := float64(fieldValue.Uint())
			metric.Value = &val
		case "uint32":
			val := float64(fieldValue.Uint())
			metric.Value = &val
		case "int64":
			val := float64(fieldValue.Int())
			metric.Value = &val
		case "float64":
			val := float64(fieldValue.Float())
			metric.Value = &val
		default:
			fmt.Printf("Unsupported value type in struct %s", fieldType)
			continue
		}
		metric.Name = key
		metric.MType = "gauge"
		metrics = append(metrics, metric)

	}
	return metrics
}

func createMetricsMap(m runtime.MemStats) []models.Metrics {

	pollCount++
	randomValue = rand.Float64() * 1000000

	metricsMap := []models.Metrics{}

	metricsMap = append(metricsMap, models.Metrics{
		Name:  "PollCount",
		MType: "counter",
		Delta: &pollCount,
	}, models.Metrics{
		Name:  "RandomValue",
		MType: "gauge",
		Value: &randomValue,
	})
	metricsMap = append(metricsMap, convertToMetrics(m)...)
	return metricsMap
}
