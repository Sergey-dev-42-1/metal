package controller

import (
	"metal/internal/pkg/domain/repositories"
	"metal/internal/server/presentation/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func testRequest(t *testing.T, gs *gin.Engine, method,
	path string) (*httptest.ResponseRecorder, string) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest(method, "http://localhost:8080"+path, nil)
	gs.ServeHTTP(w, req)

	respBody := w.Body.String()
	require.NoError(t, err)

	return w, respBody
}

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
			request: "/update/counter/test/123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := zap.SugaredLogger{}
			r := router.Router()
			ms := repositories.NewMemStorage("./save.json", &l)
			mc := New(r, &l, ms)
			r = mc.AddRoutes()
			rec, _ := testRequest(t, r, "POST", tt.request)
			assert.Equal(t, tt.want.statusCode, rec.Code)
			// assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
