package update

import (
	"fmt"
	"metal/internal/pkg/domain/models"
	"net/http"
	"net/url"
)

type UpdateService struct {
}

func (s *UpdateService) UpdateMetrics(metric models.Metric) (*http.Response, error) {

	link := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
		Path:   "/update/" + metric.Type + "/" + metric.Name + "/",
	}
	if metric.Type == "counter" {
		link.Path += metric.Value.ToString()
	} else {
		link.Path += metric.Value.ToStringFloat()
	}
	fmt.Println("Updating metrics on server", link.String())
	res, err := http.Post(link.String(), "text/plain", nil)
	fmt.Println(res)
	return res, err
}
