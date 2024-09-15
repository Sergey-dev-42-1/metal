package controller

import (
	"encoding/json"
	"fmt"
	"metal/internal/pkg/domain/models"
	"metal/internal/pkg/domain/repositories/interfaces"
	service "metal/internal/server/application/metrics-service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MetricsController struct {
	metricService *service.MetricsService
	router        *gin.Engine
}

func New(r *gin.Engine, l *zap.SugaredLogger, repo interfaces.MetricsStorage) *MetricsController {
	m := service.New(repo, l)
	return &MetricsController{
		router:        r,
		metricService: m,
	}
}

// Это было в роутере, но тогда не получится протестировать контроллер
// т.к. будет циклическая зависимость между ним и роутером (роутер -> контроллер -> роутер)
func (mc *MetricsController) AddRoutes() *gin.Engine {
	root := mc.router.Group("/")
	{
		root.GET("/", mc.HandleGetStoredValuesHTML)
		root.GET("/ping", mc.Ping)
		update := root.Group("/update")
		{
			update.POST("/", mc.HandleMetricRecordingJSON)
			update.POST(":type/:name/:value", mc.HandleMetricRecording)
		}
		value := root.Group("/value")
		{
			value.POST("/", mc.HandleGetMetricValueJSON)
			value.GET(":type/:name", mc.HandleGetMetricValue)
		}
	}
	return mc.router
}
func (mc *MetricsController) HandleMetricRecordingJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.String(415, "%s", "Request content is not marked as JSON")
		return
	}
	metric := models.Metrics{}
	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(500, "%s", "Something went wrong when trying to parse request content")
		return
	}
	fmt.Println("test record", metric.ID, metric.Delta, metric.Value, metric.MType)
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

	result := mc.metricService.CreateOrUpdateMetric(metric)
	c.JSON(200, result)
}

func (mc *MetricsController) HandleMetricRecording(c *gin.Context) {
	// fmt.Println("Metric controller", c.Params)
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
		return
	}
	metric.Value = &mv
	if v, err := strconv.ParseInt(value, 0, 64); err == nil {
		metric.Delta = &v
	}

	result := mc.metricService.CreateOrUpdateMetric(metric)
	c.String(200, "Successfully written '%s' metric", result.ID)
}

func (mc *MetricsController) HandleGetMetricValue(c *gin.Context) {

	name := c.Param("name")

	if name == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}
	metric, err := mc.metricService.FindMetric(name)
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

func (mc *MetricsController) HandleGetMetricValueJSON(c *gin.Context) {
	if c.Request.Header.Get("Content-Type") != "application/json" {
		c.String(415, "%s", "Request content is not marked as JSON")
		return
	}
	var metric models.Metrics
	if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
		c.String(500, "%s", "Something went wrong when trying to parse request content")
		return
	}

	if metric.ID == "" {
		c.String(404, "%s", "Name of the metric is not specified")
		return
	}
	m, err := mc.metricService.FindMetric(metric.ID)
	// fmt.Println("test read", metric.ID, metric.Delta, metric.Value, metric.MType)
	if err != nil {
		c.String(404, "%s", "Didn't find such metric")
		return
	}
	c.JSON(200, m)
}

func (mc *MetricsController) HandleGetStoredValuesHTML(c *gin.Context) {
	fmt.Println("Get metric page")
	metrics := mc.metricService.GetAllMetrics()

	c.HTML(200, "index.html", metrics)
}

func (mc *MetricsController) Ping(c *gin.Context) {
	fmt.Println("Get ping")
	err := mc.metricService.Ping()
	if err != nil {
		c.String(500, "%s", "Something went wrong")
		return
	}
	c.String(200, "Ping!")
}
