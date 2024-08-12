package interfaces

type MetricsService interface{
	CollectMemStats()
	SendMemStats()
}