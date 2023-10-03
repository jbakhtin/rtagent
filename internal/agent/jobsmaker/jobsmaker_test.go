package jobsmaker

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/aggregator"
	jobqueue "github.com/jbakhtin/rtagent/internal/agent/jobsqueue"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	poolCounter := aggregator.PoolCounter{Count: 0}
	agr, err := aggregator.New().WithCustomCollectors(poolCounter.PoolCount).Build()
	if err != nil {
		t.Errorf("aggregator error: %v", err)
	}

	queue := jobqueue.NewQueue()

	type args struct {
		slicer Slicer
		jober  Jober
	}
	tests := []struct {
		name string
		args args
		want *jobsMaker
	}{
		{
			"New jober",
			args{
				agr,
				queue,
			},
			&jobsMaker{
				agr,
				queue,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.slicer, tt.args.jober); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_jobsMaker_Do(t *testing.T) {
	poolCounter := aggregator.PoolCounter{Count: 0}
	agr, err := aggregator.New().WithCustomCollectors(poolCounter.PoolCount).Build()
	if err != nil {
		t.Errorf("aggregator error: %v", err)
	}

	queue := jobqueue.NewQueue()

	type fields struct {
		slicer Slicer
		jober  Jober
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
			"New jober",
			fields{
				agr,
				queue,
			},
			args{
				context.TODO(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jm := &jobsMaker{
				slicer: tt.fields.slicer,
				jober:  tt.fields.jober,
			}

			if err := jm.Do(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
