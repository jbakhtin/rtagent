package memstorage

import (
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"go.uber.org/zap"
	"reflect"
	"sync"
	"testing"
)

func TestMemStorage_Get(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metric, _ := models.NewCounter("counter", "TestMetric", "10")
	type fields struct {
		Mx     *sync.RWMutex
		Items  map[string]models.Metricer
		Logger *zap.Logger
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.Metricer
		wantErr bool
	}{
		{
			"Try to get value from empty db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{"TestMetric": metric},
				logger,
			},
			args{
				"TestMetric",
			},
			metric,
			false,
		},
		{
			"Try to get value from not empty db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{},
				logger,
			},
			args{
				"TestMetric",
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Mx:     tt.fields.Mx,
				Items:  tt.fields.Items,
				Logger: tt.fields.Logger,
			}
			got, err := ms.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetAll(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metric1, _ := models.NewCounter("counter", "TestMetric", "10")
	metric2, _ := models.NewCounter("gauge", "TestMetric2", "10")
	type fields struct {
		Mx     *sync.RWMutex
		Items  map[string]models.Metricer
		Logger *zap.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]models.Metricer
		wantErr bool
	}{
		{
			"Get all metrics from db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{"TestMetric1": metric1, "TestMetric2": metric2},
				logger,
			},
			map[string]models.Metricer{"TestMetric1": metric1, "TestMetric2": metric2},
			false,
		},
		{
			"Get all metrics from empty db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{},
				logger,
			},
			map[string]models.Metricer{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Mx:     tt.fields.Mx,
				Items:  tt.fields.Items,
				Logger: tt.fields.Logger,
			}
			got, err := ms.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_Set(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metric1, _ := models.NewCounter("counter", "TestMetric", "10")
	metric2, _ := models.NewCounter("gauge", "TestMetric2", "10")
	metric3, _ := models.NewCounter("gauge", "TestMetric2Set", "0")
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
		{
			"Get all metrics from db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{"TestMetric1": metric1, "TestMetric2": metric2},
				logger,
			},
			args{
				metric3,
			},
			metric3,
			false,
		},
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

func TestMemStorage_SetBatch(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	metric1, _ := models.NewCounter("counter", "TestMetric", "10")
	metric2, _ := models.NewCounter("gauge", "TestMetric2", "10")
	metric3, _ := models.NewCounter("gauge", "TestMetricSet1", "0")
	metric4, _ := models.NewCounter("gauge", "TestMetricSet2", "1")
	type fields struct {
		Mx     *sync.RWMutex
		Items  map[string]models.Metricer
		Logger *zap.Logger
	}
	type args struct {
		metrics []models.Metricer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Metricer
		wantErr bool
	}{
		{
			"Get all metrics from db",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{"TestMetric1": metric1, "TestMetric2": metric2},
				logger,
			},
			args{
				[]models.Metricer{metric3, metric4},
			},
			[]models.Metricer{metric3, metric4},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Mx:     tt.fields.Mx,
				Items:  tt.fields.Items,
				Logger: tt.fields.Logger,
			}
			got, err := ms.SetBatch(tt.args.metrics)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetBatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMemStorage(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	//metric1, _ := models.NewCounter("counter", "TestMetric", "10")
	//metric2, _ := models.NewCounter("gauge", "TestMetric2", "10")
	//metric3, _ := models.NewCounter("gauge", "TestMetricSet1", "0")
	//metric4, _ := models.NewCounter("gauge", "TestMetricSet2", "1")

	cfg, _ := config.NewConfigBuilder().WithAllFromEnv().Build()

	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    MemStorage
		wantErr bool
	}{
		{
			"New mem storage",
			args{
				cfg,
			},
			MemStorage{
				&sync.RWMutex{},
				map[string]models.Metricer{},
				logger,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMemStorage(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf("NewMemStorage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	type fields struct {
		Mx     *sync.RWMutex
		Items  map[string]models.Metricer
		Logger *zap.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"",
			fields{
				&sync.RWMutex{},
				map[string]models.Metricer{},
				logger,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &MemStorage{
				Mx:     tt.fields.Mx,
				Items:  tt.fields.Items,
				Logger: tt.fields.Logger,
			}
			if err := ms.TestPing(); (err != nil) != tt.wantErr {
				t.Errorf("TestPing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
