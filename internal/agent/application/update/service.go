package update

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"metal/internal/agent/application/update/interfaces"
	"metal/internal/pkg/domain/models"

	"github.com/go-resty/resty/v2"
)

type UpdateService struct {
	addr   string
	client *resty.Client
}

func New(a string) interfaces.UpdateService {
	client := resty.New()
	client.BaseURL = "http://" + a
	return &UpdateService{
		addr:   a,
		client: client,
	}
}

func (s *UpdateService) UpdateMetrics(metric models.Metrics) *resty.Response {
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
	res, _ := s.client.R().SetPathParams(p).Post("/update/{type}/{name}/{value}")
	return res
}

func compressJSON(w io.Writer, i interface{}) error {
	gz := gzip.NewWriter(w)
	if err := json.NewEncoder(gz).Encode(i); err != nil {
		return err
	}
	return gz.Close()
}
func (s *UpdateService) UpdateMetricsJSON(metric models.Metrics) *resty.Response {
	r, w := io.Pipe()
	go func() {
		err := compressJSON(w, metric)
		w.CloseWithError(err)
	}()
	fmt.Printf("Updating metrics on server %s: \n", metric.ID)
	headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}
	res, _ := s.client.R().SetHeaders(headers).SetBody(r).Post("update")
	// fmt.Println(string(res.Body()))
	return res
}
