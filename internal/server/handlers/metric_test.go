package handlers

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"github.com/jbakhtin/rtagent/internal/repositories/storages/inmemory"
	"github.com/stretchr/testify/assert"
)

func TestHandlerMetric_Get(t *testing.T) {
	type fields struct {
		repo    interfaces.MetricRepository
		request string
	}
	type want struct {
		statusCode  int
	}

	storage, err := inmemory.NewMetricRepository()
	if err != nil {
		fmt.Println(err)
		return
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			"Get all metrics. Test 1",
			fields{
				repo:    storage,
				request: "http://127.0.0.1:8080",
			},
			want{
				statusCode:  200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerMetric{
				repo: tt.fields.repo,
			}

			request := httptest.NewRequest(http.MethodGet, tt.fields.request, nil)
			w := httptest.NewRecorder()
			hf := h.GetAll()
			hf(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}