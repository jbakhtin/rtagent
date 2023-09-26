package aggregator

import (
	"context"
	"github.com/go-faster/errors"
	"github.com/jbakhtin/rtagent/internal/types"
	gopsutil "github.com/shirou/gopsutil/v3/mem"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestBuilder_Build(t *testing.T) {
	aggr, err := New().WithDefaultCollectors().Build()
	if err != nil {
		t.Errorf("Build() error = %v", err)
	}

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
		{
			"Test Build with out errors",
			fields{
				nil,
				aggregator{
					[]CollectorFunc{Runtime, Gopsutil, RandomMetric},
					Metrics{
						make(map[string]types.Metricer),
						sync.RWMutex{},
					},
					sync.RWMutex{},
				},
			},
			aggr,
			false,
		},
		{
			"Test Build with errors",
			fields{
				errors.New("Some error in builder"),
				aggregator{
					[]CollectorFunc{Runtime, Gopsutil, RandomMetric},
					Metrics{
						make(map[string]types.Metricer),
						sync.RWMutex{},
					},
					sync.RWMutex{},
				},
			},
			aggr,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			got, err := b.Build()
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				return
			}

			if len(got.collectors) != len(tt.want.collectors) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got.collection, tt.want.collection) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got.RWMutex, tt.want.RWMutex) {
				t.Errorf("Build() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestBuilder_WithCustomCollector(t *testing.T) {
	aggrBuilder := New().WithCustomCollector(RandomMetric)

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
		{
			"Test Builder with custom collector",
			fields{
				nil,
				aggregator{},
			},
			args{
				RandomMetric,
			},
			aggrBuilder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			got := b.WithCustomCollectors(tt.args.collector)
			if len(got.aggregator.collectors) != len(tt.want.aggregator.collectors) {
				t.Errorf("WithDefaultCollector() = %v, want %v", got.aggregator.collectors, tt.want.aggregator.collectors)
				return
			}
		})
	}
}

func TestBuilder_WithCustomCollectors(t *testing.T) {
	aggrBuilder := New().WithCustomCollectors(RandomMetric)


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
		{
			"Test aggregator with default collectors",
			fields{
				nil,
				aggregator{},
			},
			args{
				[]CollectorFunc{RandomMetric},
			},
			aggrBuilder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}
			got := b.WithCustomCollectors(tt.args.collectors...)
			if len(got.aggregator.collectors) != len(tt.want.aggregator.collectors) {
				t.Errorf("WithDefaultCollectors() = %v, want %v", got.aggregator.collectors, tt.want.aggregator.collectors)
				return
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
		{
			"Test aggregator with default collectors",
			fields{
				nil,
				aggregator{},
			},
			&Builder{
				nil,
				aggregator{
					collectors: []CollectorFunc{Runtime, Gopsutil, RandomMetric},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Builder{
				err:        tt.fields.err,
				aggregator: tt.fields.aggregator,
			}

			got := b.WithDefaultCollectors()

			if len(got.aggregator.collectors) != len(tt.want.aggregator.collectors) {
				t.Errorf("WithDefaultCollectors() = %v, want %v", got.aggregator.collectors, tt.want.aggregator.collectors)
				return
			}
			//for i := range got.aggregator.collectors {
			//	if !reflect.DeepEqual(got.aggregator.collectors[i], tt.want.aggregator.collectors[i]) {
			//		t.Errorf("WithDefaultCollectors() = %v, want %v", got.aggregator.collectors, tt.want.aggregator.collectors)
			//	}
			//}
		})
	}
}

func TestGopsutil(t *testing.T) {
	memStats, err := gopsutil.VirtualMemory()
	if err != nil {
		t.Errorf("Gopsutil() error = %v", err)
		return
	}
	tests := []struct {
		name    string
		want    map[string]types.Metricer
		wantErr bool
	}{
		{
			"Get Gopsutil",
			map[string]types.Metricer{
				"TotalMemory": types.Gauge(memStats.Total),
				"FreeMemory": types.Gauge(memStats.Free),
				"CPUutilization1": types.Gauge(memStats.Used),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Gopsutil()
			if (err != nil) != tt.wantErr {
				t.Errorf("Gopsutil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) !=  len(tt.want) {
				t.Errorf("Gopsutil() len(got) = %v, len(want) %v", len(got), len(tt.want))
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
		{
			"Get all Metrics empty collection",
			fields{
				map[string]types.Metricer{},
				sync.RWMutex{},
			},
			map[string]types.Metricer{},
		},
		{
			"Get all Metrics not empty collection",
			fields{
				map[string]types.Metricer{"TestMetric": types.Gauge(10)},
				sync.RWMutex{},
			},
			map[string]types.Metricer{"TestMetric": types.Gauge(10)},
		},
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
		{
			"Merge two Metrics map",
			fields{
				map[string]types.Metricer{"PoolCount" : types.Counter(10), "RandomMetric" : types.Gauge(10)},
				sync.RWMutex{},
			},
			args{
				map[string]types.Metricer{"TestCount" : types.Counter(20)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collection := &Metrics{
				items:   tt.fields.items,
				RWMutex: tt.fields.RWMutex,
			}

			collection.Merge(tt.args.items)

			itemsCount := len(collection.items)

			if itemsCount != 3 {
				t.Errorf("count items after Merge: %v, want %v", itemsCount, 3)
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
		{
			"New builder",
			&Builder{
				aggregator: aggregator{
					collection: Metrics{
						items: make(map[string]types.Metricer, 0),
					},
				},
			},
		},
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
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	tests := []struct {
		name    string
		want    map[string]types.Metricer
		wantErr bool
	}{
		{
			"Get Runtime",
			map[string]types.Metricer{
				"Alloc": types.Gauge(memStats.Alloc),
				"Frees": types.Gauge(memStats.Frees),
				"HeapAlloc": types.Gauge(memStats.HeapAlloc),
				"BuckHashSys": types.Gauge(memStats.BuckHashSys),
				"GCSys": types.Gauge(memStats.GCSys),
				"GCCPUFraction": types.Gauge(memStats.GCCPUFraction),
				"HeapIdle": types.Gauge(memStats.HeapIdle),
				"HeapInuse": types.Gauge(memStats.HeapInuse),
				"HeapObjects": types.Gauge(memStats.HeapObjects),
				"HeapReleased": types.Gauge(memStats.HeapReleased),
				"HeapSys": types.Gauge(memStats.HeapSys),
				"LastGC": types.Gauge(memStats.LastGC),
				"Lookups": types.Gauge(memStats.Lookups),
				"MCacheInuse": types.Gauge(memStats.MCacheInuse),
				"MCacheSys": types.Gauge(memStats.MCacheSys),
				"MSpanInuse": types.Gauge(memStats.MSpanInuse),
				"MSpanSys": types.Gauge(memStats.MSpanSys),
				"Mallocs": types.Gauge(memStats.Mallocs),
				"NextGC": types.Gauge(memStats.NextGC),
				"NumForcedGC": types.Gauge(memStats.NumForcedGC),
				"NumGC": types.Gauge(memStats.NumGC),
				"OtherSys": types.Gauge(memStats.OtherSys),
				"PauseTotalNs": types.Gauge(memStats.PauseTotalNs),
				"StackInuse": types.Gauge(memStats.StackInuse),
				"StackSys": types.Gauge(memStats.StackSys),
				"Sys": types.Gauge(memStats.Sys),
				"TotalAlloc": types.Gauge(memStats.TotalAlloc),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Runtime()
			if (err != nil) != tt.wantErr {
				t.Errorf("Runtime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) !=  len(tt.want) {
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
		{
			"Get All",
			fields{
				[]CollectorFunc{Gopsutil},
				Metrics{
					map[string]types.Metricer{"TestMetric": types.Gauge(10)},
					sync.RWMutex{},
				},
				sync.RWMutex{},
				1,
			},
			map[string]types.Metricer{"TestMetric": types.Gauge(10)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				collectors: tt.fields.collectors,
				collection: tt.fields.collection,
				RWMutex:    tt.fields.RWMutex,
			}
			if got := a.GetAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_aggregator_Pool(t *testing.T) {
	pooCounter := PoolCounter{}
	type fields struct {
		collectors []CollectorFunc
		collection Metrics
		RWMutex    sync.RWMutex
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
		{
			name: "Aggregator with counter func",
			fields: fields{
				[]CollectorFunc{pooCounter.PoolCount},
				Metrics{
					items: make(map[string]types.Metricer),
				},
				sync.RWMutex{},
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				collectors: tt.fields.collectors,
				collection: tt.fields.collection,
				RWMutex:    tt.fields.RWMutex,
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
			poolCounter := PoolCounter{}
			got, err := poolCounter.PoolCount()
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
