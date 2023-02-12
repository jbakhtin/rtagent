package handlers

import (
	"github.com/jbakhtin/rtagent/internal/repositories/interfaces"
	"net/http"
	"reflect"
	"testing"
)

func TestHandlerMetric_Find(t *testing.T) {
	type fields struct {
		repo interfaces.MetricRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HandlerMetric{
				repo: tt.fields.repo,
			}
			if got := h.Find(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}