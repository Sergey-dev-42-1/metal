package controller

import (
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
			update.POST(":type/:name/:value", HandleMetricRecording)
		}
		value := root.Group("/value/")
		{
			value.GET(":type/:name", HandleGetMetricValue)
		}
	}
	return mc.r
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

	metricValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		c.String(400, "%s", "Value is not a number")
		return
	}

	result := service.CreateOrUpdateMetric(models.Metric{Value: models.MetricValue(metricValue), Type: tp, Name: name})
	c.String(200, "Successfully written '%s' metric", result.Name)
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
	if metric.Type == "counter" {
		c.String(200, "%s", metric.Value.ToString())
		return
	}
	c.String(200, "%s", metric.Value.ToStringFloat())
}
func HandleGetStoredValuesHTML(c *gin.Context) {
	fmt.Println("Get metric page")
	viewData := service.GetAllMetrics()
	c.HTML(200, "index.html", viewData)
}
