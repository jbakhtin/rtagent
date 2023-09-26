package once

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/agent/tasker"
	"reflect"
	"sync"
	"testing"
	"time"
)

type counter struct {
	Counter int
	sync.RWMutex
}

func (c *counter) Increment(ctx context.Context) error {
	c.Lock()
	defer c.Unlock()
	c.Counter++
	return nil
}

func (c *counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.Counter
}

func TestNew(t *testing.T) {
	counter := counter{}
	type args struct {
		name string
		f    tasker.Func
	}
	tests := []struct {
		name string
		args args
		want *task
	}{
		{
			"New Once task",
			args{
				"Once task",
				counter.Increment,
			},
			&task{
				counter.Increment,
				"Once task",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.name, tt.args.f)

			if !reflect.DeepEqual(got.name, tt.want.name) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_task_Do(t1 *testing.T) {
	counter := counter{}
	type fields struct {
		f    tasker.Func
		name string
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
			"New Once task",
			fields{
				counter.Increment,
				"Once task",
			},
			args{
				context.TODO(),
			},
			false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &task{
				f:    tt.fields.f,
				name: tt.fields.name,
			}
			timer := time.NewTimer(time.Second * 1)
			go func() {
				if err := t.Do(tt.args.ctx); (err != nil) != tt.wantErr {
					t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			<-timer.C

			if counter.Get() != 1 {
				t1.Errorf("A different result was expected: %v, want: %v", counter.Counter, 1)
			}
		})
	}
}
