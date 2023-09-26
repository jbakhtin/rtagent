package aggregator

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/types"
	"reflect"
	"sync"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	type fields struct {
		err        error
		aggregator aggregator
	}
	tests := []struct {
		name    string
		fields  fields
		want    *aggregator
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			got, err := b.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithCustomCollector(t *testing.T) {
	type fields struct {
		err        error
		aggregator aggregator
	}
	type args struct {
		collector CollectorFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			if got := b.WithCustomCollector(tt.args.collector); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCustomCollector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithCustomCollectors(t *testing.T) {
	type fields struct {
		err        error
		aggregator aggregator
	}
	type args struct {
		collectors []CollectorFunc
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			if got := b.WithCustomCollectors(tt.args.collectors); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCustomCollectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder_WithDefaultCollectors(t *testing.T) {
	type fields struct {
		err        error
		aggregator aggregator
	}
	tests := []struct {
		name   string
		fields fields
		want   *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			if got := b.WithDefaultCollectors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDefaultCollectors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGopsutil(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]types.Metricer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gopsutil()
			if (err != nil) != tt.wantErr {
				t.Errorf("Gopsutil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gopsutil() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_GetAll(t *testing.T) {
	type fields struct {
		items   map[string]types.Metricer
		RWMutex sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]types.Metricer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				items:   tt.fields.items,
				RWMutex: tt.fields.RWMutex,
			}
			if got := m.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetrics_Merge(t *testing.T) {
	type fields struct {
		items   map[string]types.Metricer
		RWMutex sync.RWMutex
	}
	type args struct {
		items map[string]types.Metricer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Metrics{
				items:   tt.fields.items,
				RWMutex: tt.fields.RWMutex,
			}
		})
	}
}

func TestMetrics_Set(t *testing.T) {
	type fields struct {
		items   map[string]types.Metricer
		RWMutex sync.RWMutex
	}
	type args struct {
		key    string
		metric types.Metricer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Set metric",
			fields: fields{
				map[string]types.Metricer{"PoolCount" : types.Counter(1)},
				sync.RWMutex{},
			},
			args: args{
				"PoolCount",
				types.Counter(1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Metrics{
				items:   tt.fields.items,
				RWMutex: tt.fields.RWMutex,
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Builder
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRandomMetric(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]types.Metricer
		wantErr bool
	}{
		{
			"Get a random metric",
			map[string]types.Metricer{"RandomMetric":nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomMetric()
			if (err != nil) != tt.wantErr {
				t.Errorf("RandomMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_, ok := tt.want["RandomMetric"]
			_, ok2 := got["RandomMetric"]

			if ok && ok2 {
				t.Errorf("RandomMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRuntime(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]types.Metricer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Runtime()
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Runtime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregator_GetAll(t *testing.T) {
	type fields struct {
		collectors []CollectorFunc
		collection Metrics
		RWMutex    sync.RWMutex
		poolCount  types.Counter
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]types.Metricer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				collectors: tt.fields.collectors,
				collection: tt.fields.collection,
				RWMutex:    tt.fields.RWMutex,
				poolCount:  tt.fields.poolCount,
			}
			if got := a.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregator_Pool(t *testing.T) {
	type fields struct {
		collectors []CollectorFunc
		collection Metrics
		RWMutex    sync.RWMutex
		poolCount  types.Counter
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				collectors: tt.fields.collectors,
				collection: tt.fields.collection,
				RWMutex:    tt.fields.RWMutex,
				poolCount:  tt.fields.poolCount,
			}
			if err := a.Pool(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Pool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_aggregator_poolCountCollector(t *testing.T) {
	type fields struct {
		collectors []CollectorFunc
		collection Metrics
		RWMutex    sync.RWMutex
		poolCount  types.Counter
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]types.Metricer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				collectors: tt.fields.collectors,
				collection: tt.fields.collection,
				RWMutex:    tt.fields.RWMutex,
				poolCount:  tt.fields.poolCount,
			}
			got, err := a.poolCountCollector()
			if (err != nil) != tt.wantErr {
				t.Errorf("poolCountCollector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("poolCountCollector() got = %v, want %v", got, tt.want)
			}
		})
	}
}
