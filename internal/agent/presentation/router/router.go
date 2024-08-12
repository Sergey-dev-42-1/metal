package router

import (
	"metal/internal/server/presentation/controller"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/update/{type}/{name}/{value}", controller.HandleMetricRecording)
	mux.HandleFunc("/find/{name}", controller.HandleFindMetric)
	// mux.HandleFunc("/update/{type}/{name}/{value}", controller.HandleNotExisting)
	return mux
}
