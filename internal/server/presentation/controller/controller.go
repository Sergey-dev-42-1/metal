package controller

import (
	"encoding/json"
	"fmt"
	"metal/internal/pkg/domain/models"
	"metal/internal/server/application/metrics-service"
	"net/http"
	"strconv"
)

func HandleMetricRecording(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Metric controller", req.URL.Path)
	value := req.PathValue("value")
	tp := req.PathValue("type")
	name := req.PathValue("name")
	
	if name == "" {
		http.Error(res, "Name of the metric is not specified", http.StatusNotFound)
		return
	}

	if (tp != "gauge" && tp != "counter") || value == "" {
		http.Error(res, "Bad request, check parameters", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case "POST":
		{
			
			metricValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				http.Error(res, "Value is not a number", http.StatusBadRequest)
			}
			
			service.CreateOrUpdateMetric(models.Metric{Value: models.MetricValue(metricValue), Type: tp, Name: name})
			// jsonMetric, _ := json.Marshal(metric)
			// res.Write(jsonMetric)
			return
		}
	default:
		{
			http.Error(res, "Method not supported", http.StatusMethodNotAllowed)
			return
		}
	}

}

func HandleFindMetric(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Find metric")
	name := req.PathValue("name")

	if name == "" {
		http.Error(res, "Name of the metric is not specified", http.StatusNotFound)
	}
	switch req.Method {
	case "GET":
		{

			metricValue, err := service.FindMetric(name)
			if err != nil {
				http.Error(res, "Error while getting metric", http.StatusInternalServerError)
			}
			jsonMetric, _ := json.Marshal(metricValue)
			res.Write(jsonMetric)
			res.Write([]byte(name + " " + string(jsonMetric)))
		}
	default:
		{
			http.Error(res, "Method not supported", http.StatusMethodNotAllowed)
		}
	}

}
