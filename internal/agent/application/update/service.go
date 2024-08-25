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

func (s *UpdateService) UpdateMetrics(metric models.Metric) (*resty.Response, error) {

	client := resty.New()
	client.BaseURL = "http://" + s.addr
	p := map[string]string{
		"type": metric.Type,
		"name": metric.Name,
	}

	if metric.Type == "counter" {
		p["value"] += metric.Value.ToString()
	} else {
		p["value"] = metric.Value.ToStringFloat()
	}

	fmt.Printf("Updating metrics on server %s:%f \n", metric.Name, metric.Value)
	res, err := client.R().SetPathParams(p).Post("/update/{type}/{name}/{value}")

	if err != nil {
		panic(err)
	}
	return res, err
}
