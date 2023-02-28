package handlers

import (
	"context"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/stretchr/testify/require"
	"log"
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

	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()
	storage, err := inmemory.NewMetricRepository(ctx, cfg)
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
			var cfg config.Config
			err := env.Parse(&cfg)
			if err != nil {
				log.Println(err)
			}

			ctx := context.Background()
			h, err := NewHandlerMetric(ctx, cfg)
			if err != nil {
				require.NoError(t, err)
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
