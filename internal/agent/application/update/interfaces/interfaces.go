package interfaces

import (
	"metal/internal/pkg/domain/models"
	"net/http"
)

type UpdateService interface {
	UpdateMetrics(models.Metric) (http.Response, error)
}
