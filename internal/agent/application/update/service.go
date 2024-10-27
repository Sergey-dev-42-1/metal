package update

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"metal/internal/agent/application/update/interfaces"

	// "metal/internal/pkg/crypto"

	"metal/internal/pkg/crypto"
	"metal/internal/pkg/domain/models"

	"github.com/go-resty/resty/v2"
)

type UpdateService struct {
	addr   string
	key    string
	client *resty.Client
}

func New(a string, k string) interfaces.UpdateService {
	client := resty.New()
	client.BaseURL = "http://" + a
	return &UpdateService{
		addr:   a,
		key:    k,
		client: client,
	}
}

func (s *UpdateService) UpdateMetrics(metric models.Metrics) *resty.Response {
	p := map[string]string{
		"type": metric.MType,
		"name": metric.Name,
	}

	if metric.MType == "counter" {
		m := *metric.Delta
		p["value"] = fmt.Sprintf("%d", m)
	} else {
		m := *metric.Value
		p["value"] = fmt.Sprintf("%g", m)
	}

	fmt.Printf("Updating metrics on server %s:%f \n", metric.Name, *metric.Value)
	res, _ := s.client.R().SetPathParams(p).Post("/update/{type}/{name}/{value}")
	return res
}

func compressJSON(w io.Writer, jsonData []byte) error {
	gzipWriter := gzip.NewWriter(w)
	_, err := gzipWriter.Write(jsonData)
	if err != nil {
		log.Fatal("Error writing compressed data:", err)
	}
	gzipWriter.Close()
	return err
}

// TODO: could just use batch one but pass array with one element
func (s *UpdateService) UpdateMetricsJSON(metric models.Metrics) *resty.Response {
	r, w := io.Pipe()
	jsonData, err := json.Marshal(metric)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}
	go func() {
		compressErr := compressJSON(w, jsonData)
		w.CloseWithError(compressErr)
	}()
	fmt.Printf("Updating metrics on server %s: \n", metric.Name)
	headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}
	if s.key != "" {
		hash := crypto.SignSHA256(jsonData, s.key)
		headers["HashSHA256"] = hash
	}

	res, err := s.client.R().SetHeaders(headers).SetBody(r).Post("update")
	fmt.Println(res, err)
	return res
}

func (s *UpdateService) UpdateMetricsJSONBatch(metric []models.Metrics) *resty.Response {
	r, w := io.Pipe()
	jsonData, err := json.Marshal(metric)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}
	go func() {
		compressErr := compressJSON(w, jsonData)
		w.CloseWithError(compressErr)
	}()
	fmt.Printf("Updating metrics on server \n")
	headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip"}
	if s.key != "" {
		hash := string(crypto.SignSHA256(jsonData, s.key))
		headers["HashSHA256"] = hash
	}
	res, _ := s.client.R().SetHeaders(headers).SetBody(r).Post("updates")
	// fmt.Println(string(res.Body()))
	return res
}
