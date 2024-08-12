package service

import (
	"fmt"
	"math/rand"
	"metal/internal/agent/application/update"
	"metal/internal/pkg/domain/models"
	"runtime"
	"time"
)

var stats []models.Metric
var service = update.UpdateService{}
var pollCount int64
var randomValue int64

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
func createMetricsMap(m runtime.MemStats) []models.Metric {

	pollCount++
	randomValue = rand.Int63n(1000000)

	metricsMap := []models.Metric{}
	//как-то рефлексией это решить?
	metricsMap = append(metricsMap, models.Metric{
		Value: models.MetricValue(pollCount),
		Type:  "counter",
		Name:  "PollCount",
	}, models.Metric{
		Value: models.MetricValue(randomValue),
		Type:  "gauge",
		Name:  "RandomValue",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.Alloc),
		Name:  "Alloc",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.BuckHashSys),
		Name:  "BuckHashSys",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.Frees),
		Name:  "Frees",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.GCCPUFraction),
		Name:  "GCCPUFraction",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.GCSys),
		Name:  "GCSys",
	}, models.Metric{
		Type:  "gauge",
		Value: models.MetricValue(m.HeapAlloc),
		Name:  "HeapAlloc",
	},
	)

	return metricsMap
	// metricsMap["HeapIdle"] = models.Metric{
	// 	Value: models.MetricValue(m.HeapIdle),
	// }
	// metricsMap["HeapInuse"] = models.Metric{
	// 	Value: models.MetricValue(m.HeapInuse),
	// }
	// metricsMap["HeapObjects"] = models.Metric{
	// 	Value: models.MetricValue(m.HeapObjects),
	// }
	// metricsMap["HeapReleased"] = models.Metric{
	// 	Value: models.MetricValue(m.HeapReleased),
	// }
	// metricsMap["HeapSys"] = models.Metric{
	// 	Value: models.MetricValue(m.HeapSys),
	// }
	// metricsMap["LastGC"] = models.Metric{
	// 	Value: models.MetricValue(m.LastGC),
	// }
	// metricsMap["Lookups"] = models.Metric{
	// 	Value: models.MetricValue(m.Lookups),
	// }
	// metricsMap["MCacheInuse"] = models.Metric{
	// 	Value: models.MetricValue(m.MCacheInuse),
	// }
	// metricsMap["MCacheSys"] = models.Metric{
	// 	Value: models.MetricValue(m.MCacheSys),
	// }
	// metricsMap["MSpanInuse"] = models.Metric{
	// 	Value: models.MetricValue(m.MSpanInuse),
	// }
	// metricsMap["MSpanSys"] = models.Metric{
	// 	Value: models.MetricValue(m.MSpanSys),
	// }
	// metricsMap["Mallocs"] = models.Metric{
	// 	Value: models.MetricValue(m.Mallocs),
	// }
	// metricsMap["NextGC"] = models.Metric{
	// 	Value: models.MetricValue(m.NextGC),
	// }
	// metricsMap["NumForcedGC"] = models.Metric{
	// 	Value: models.MetricValue(m.NumForcedGC),
	// }
	// metricsMap["NumGC"] = models.Metric{
	// 	Value: models.MetricValue(m.NumGC),
	// }
	// metricsMap["OtherSys"] = models.Metric{
	// 	Value: models.MetricValue(m.OtherSys),
	// }
	// metricsMap["PauseTotalNs"] = models.Metric{
	// 	Value: models.MetricValue(m.PauseTotalNs),
	// }
	// metricsMap["StackInuse"] = models.Metric{
	// 	Value: models.MetricValue(m.StackInuse),
	// }
	// metricsMap["StackSys"] = models.Metric{
	// 	Value: models.MetricValue(m.StackSys),
	// }
	// metricsMap["Sys"] = models.Metric{
	// 	Value: models.MetricValue(m.Sys),
	// }
	// metricsMap["TotalAlloc"] = models.Metric{
	// 	Value: models.MetricValue(m.TotalAlloc),
	// }
}
