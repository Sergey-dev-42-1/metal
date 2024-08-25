package interfaces

import (
	"metal/internal/pkg/domain/models"
	"github.com/go-resty/resty/v2"
)

type UpdateService interface {
	UpdateMetrics(models.Metric) (*resty.Response, error)
}
