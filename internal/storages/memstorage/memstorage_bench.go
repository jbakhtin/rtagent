package memstorage

import (
	"github.com/jbakhtin/rtagent/internal/models"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_Set(t *testing.T) {
	type fields struct {
		Mx     *sync.RWMutex
		Items  map[string]models.Metricer
		Logger *zap.Logger
	}
	type args struct {
		metric models.Metricer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Metricer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Mx:     tt.fields.Mx,
				Items:  tt.fields.Items,
				Logger: tt.fields.Logger,
			}
			got, err := ms.Set(tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set() got = %v, want %v", got, tt.want)
			}
		})
	}
}
