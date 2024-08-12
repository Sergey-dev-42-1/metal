package controller

import (
	"metal/internal/pkg/domain/repositories"
	service "metal/internal/server/application/metrics-service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMetricRecording(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}
	tests := []struct {
		name    string
		request string
		pv      map[string]string
		want    want
	}{
		{
			name: "receives counter metric",
			want: want{
				statusCode: 200,
			},
			pv:      map[string]string{"type": "counter", "value": "123", "name": "test"},
			request: "/update/counter/test/123/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service.SetStorage(repositories.New())
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			for key, v := range tt.pv {
				request.SetPathValue(key, v)
			}
			w := httptest.NewRecorder()
			h := http.HandlerFunc(HandleMetricRecording)
			h(w, request)

			result := w.Result()
			result.Body.Close()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			// assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

		})
	}
}
