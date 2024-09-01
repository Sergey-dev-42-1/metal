package controller

import (
	"encoding/json"
	"fmt"
	"metal/internal/pkg/domain/models"
	service "metal/internal/server/application/metrics-service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MetricsController struct {
	r *gin.Engine
}

func New(r *gin.Engine) *MetricsController {
	return &MetricsController{
		r,
	}
}

// Это было в роутере, но тогда не получится протестировать контроллер
// т.к. будет циклическая зависимость между ним и роутером (роутер -> контроллер -> роутер)
func (mc *MetricsController) AddRoutes() *gin.Engine {
	root := mc.r.Group("/")
	{
		root.GET("/", HandleGetStoredValuesHTML)

		update := root.Group("/update/")
		{
			update.POST("", HandleMetricRecordingJSON)
			update.POST(":type/:name/:value", HandleMetricRecording)
		}
		value := root.Group("/value/")
		{
			value.POST("", HandleGetMetricValueJSON)
			value.GET(":type/:name", HandleGetMetricValue)
		}
	}
	return mc.r
}
func HandleMetricRecordingJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.String(415, "%s", "Request content is not marked as JSON")
		return
	}
	metric := models.Metrics{}
	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(500, "%s", "Something went wrong when trying to parse request content")
		return
	}

	if metric.ID == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}
	if metric.MType != "gauge" && metric.MType != "counter" {
		c.String(400, "%s", "Bad request, check parameters")
		return
	}
	if metric.MType == "gauge" && metric.Value == nil {

		c.String(400, "%s", "Bad request, check parameters")
		return
	}
	if metric.MType == "counter" && metric.Delta == nil {

		c.String(400, "%s", "Bad request, check parameters")
		return
	}

	result := service.CreateOrUpdateMetric(metric)
	c.String(200, "Successfully written '%s' metric", result.ID)
}

func HandleMetricRecording(c *gin.Context) {
	fmt.Println("Metric controller", c.Params)
	value := c.Param("value")
	tp := c.Param("type")
	name := c.Param("name")

	if name == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}

	if (tp != "gauge" && tp != "counter") || value == "" {
		c.String(400, "%s", "Bad request, check parameters")
		return
	}

	metric := models.Metrics{ID: name, MType: tp}
	mv, err := strconv.ParseFloat(value, 64)
	if err != nil {
		c.String(400, "%s", "Value is not a number")
	}
	metric.Value = &mv
	if v, err := strconv.ParseInt(value, 0, 64); err == nil {
		metric.Delta = &v
	}

	result := service.CreateOrUpdateMetric(metric)
	c.String(200, "Successfully written '%s' metric", result.ID)
}

func HandleGetMetricValue(c *gin.Context) {
	fmt.Println("Get metric value")
	name := c.Param("name")
	fmt.Println(name)
	if name == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}
	metric, err := service.FindMetric(name)
	if err != nil {
		c.String(404, "%s", "Didn't find such metric")
		return
	}
	if metric.MType == "counter" {
		c.String(200, "%s", metric.Value)
		return
	}
	c.String(200, "%s", fmt.Sprintf("%g", *metric.Value))
}

func HandleGetMetricValueJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.String(415, "%s", "Request content is not marked as JSON")
		return
	}
	metric := models.Metrics{}
	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(500, "%s", "Something went wrong when trying to parse request content")
		return
	}
	
	if metric.ID == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}
	metric, err := service.FindMetric(metric.ID)
	if err != nil {
		c.String(404, "%s", "Didn't find such metric")
		return
	}
	if metric.MType == "counter" {
		c.String(200, "%s", fmt.Sprintf("%d", *metric.Delta))
		return
	}
	c.String(200, "%s", fmt.Sprintf("%g", *metric.Value))
}

func HandleGetStoredValuesHTML(c *gin.Context) {
	fmt.Println("Get metric page")
	metrics := service.GetAllMetrics()

	c.HTML(200, "index.html", metrics)
}
