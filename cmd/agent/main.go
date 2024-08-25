package main

import service "metal/internal/agent/application/metrics"

func main() {
	parseFlags()
	metricsService := service.New(startAddress, reportInterval, pollInterval)
	go metricsService.CollectMemStats()
	go metricsService.SendMemStats()
	select {}
}
