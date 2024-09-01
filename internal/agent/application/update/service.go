package update

import (
	"fmt"
	"metal/internal/agent/application/update/interfaces"
	"metal/internal/pkg/domain/models"

	"github.com/go-resty/resty/v2"
)

type UpdateService struct {
	addr string
}

func New(a string) interfaces.UpdateService {
	return &UpdateService{
		addr: a,
	}
}

func (s *UpdateService) UpdateMetrics(metric models.Metrics) (*resty.Response) {

	client := resty.New()
	client.BaseURL = "http://" + s.addr
	p := map[string]string{
		"type": metric.MType,
		"name": metric.ID,
	}

	if metric.MType == "counter" {
		m := *metric.Delta
		p["value"] = fmt.Sprintf("%d", m)
	} else {
		m := *metric.Value
		p["value"] = fmt.Sprintf("%g", m)
	}

	fmt.Printf("Updating metrics on server %s:%f \n", metric.ID, *metric.Value)
	res, _ := client.R().SetPathParams(p).Post("/update/{type}/{name}/{value}")
	return res
}
func (s *UpdateService) UpdateMetricsJSON(metric models.Metrics) (*resty.Response) {

	client := resty.New()
	client.BaseURL = "http://" + s.addr
	fmt.Printf("Updating metrics on server %s: \n", metric.ID)
	res, _ := client.R().SetBody(&metric).Post("update")
	fmt.Println(string(res.Body()))
	return res
}
