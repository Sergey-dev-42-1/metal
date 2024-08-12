package service

import (
	"fmt"
	"math/rand"
	"metal/internal/agent/application/update"
	"metal/internal/pkg/domain/models"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"time"
)

var stats []models.Metric
var service = update.UpdateService{}
var pollCount int64
var randomValue int64
var collectedMetrics = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle",
	"HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys",
	"Sys", "TotalAlloc"}

func CollectMemStats() {
	for {
		fmt.Println("Updating metrics")
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		collectedMetrics := createMetricsMap(memStats)
		stats = collectedMetrics
		time.Sleep(2 * time.Second)
	}
}

func SendMemStats() {
	for {
		time.Sleep(10 * time.Second)
		pollCount = 0
		fmt.Println("Saving metrics")

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

	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.Alloc),
	// 	Name:  "Alloc",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.BuckHashSys),
	// 	Name:  "BuckHashSys",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.Frees),
	// 	Name:  "Frees",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.GCCPUFraction),
	// 	Name:  "GCCPUFraction",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.GCSys),
	// 	Name:  "GCSys",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapAlloc),
	// 	Name:  "HeapAlloc",
	// }, models.Metric{
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapIdle),
	// 	Name:  "HeapIdle",
	// }, models.Metric{
	// 	Name:  "HeapInuse",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapInuse),
	// }, models.Metric{
	// 	Name:  "HeapObjects",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapObjects),
	// }, models.Metric{
	// 	Name:  "HeapReleased",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapReleased),
	// }, models.Metric{
	// 	Name:  "HeapSys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.HeapSys),
	// }, models.Metric{
	// 	Name:  "LastGC",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.LastGC),
	// }, models.Metric{
	// 	Name:  "Lookups",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.Lookups),
	// }, models.Metric{
	// 	Name:  "MCacheInuse",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.MCacheInuse),
	// }, models.Metric{
	// 	Name:  "MCacheSys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.MCacheSys),
	// }, models.Metric{
	// 	Name:  "MSpanInuse",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.MSpanInuse),
	// }, models.Metric{
	// 	Name:  "MSpanSys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.MSpanSys),
	// }, models.Metric{
	// 	Name:  "Mallocs",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.Mallocs),
	// }, models.Metric{
	// 	Name:  "NextGC",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.NextGC),
	// }, models.Metric{
	// 	Name:  "NumForcedGC",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.NumForcedGC),
	// }, models.Metric{
	// 	Name:  "NumGC",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.NumGC),
	// }, models.Metric{
	// 	Name:  "OtherSys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.OtherSys),
	// }, models.Metric{
	// 	Name:  "PauseTotalNs",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.PauseTotalNs),
	// }, models.Metric{
	// 	Name:  "StackInuse",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.StackInuse),
	// }, models.Metric{
	// 	Name:  "StackSys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.StackSys),
	// }, models.Metric{
	// 	Name:  "Sys",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.Sys),
	// }, models.Metric{
	// 	Name:  "TotalAlloc",
	// 	Type:  "gauge",
	// 	Value: models.MetricValue(m.TotalAlloc),
	// },
	// )

}
