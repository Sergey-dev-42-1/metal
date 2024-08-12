package main

import service "metal/internal/agent/application/metrics"

func main() {
	go service.CollectMemStats()
	go service.SendMemStats()
	select {}
}
