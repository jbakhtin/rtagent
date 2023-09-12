package handlers

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerMetric_Get(t *testing.T) {
	type fields struct {
		request string
	}
	type want struct {
		statusCode int
	}

	ctx := context.TODO()
	cfg, err := config.NewConfigBuilder().
		WithAllFromFlagsS().
		WithAllFromEnv().
		Build()
	cfg.Restore = false

	if err != nil {
		log.Println(err)
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			"Get all metrics. Test 1",
			fields{
				request: "http://127.0.0.1:8080",
			},
			want{
				statusCode: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewHandlerMetric(ctx, cfg)
			if err != nil {
				require.NoError(t, err)
			}

			request := httptest.NewRequest(http.MethodGet, tt.fields.request, nil)
			w := httptest.NewRecorder()
			hf := h.GetAllMetricsAsHTML()
			hf(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}
