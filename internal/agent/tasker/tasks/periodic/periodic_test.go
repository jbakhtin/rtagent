package periodic

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
		name     string
		duration time.Duration
		f        tasker.Func
	}
	tests := []struct {
		name string
		args args
		want *task
	}{
		{
			"New limited task",
			args{
				"increment counter",
				time.Second * 10,
				counter.Increment,
			},
			&task{
				counter.Increment,
				"increment counter",
				time.Second * 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.name, tt.args.duration, tt.args.f)

			if !reflect.DeepEqual(got.name, tt.want.name) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(got.duration, tt.want.duration) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_task_Do(t1 *testing.T) {
	counter := counter{}
	type fields struct {
		f        tasker.Func
		name     string
		duration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Test limited task",
			fields{
				counter.Increment,
				"increment counter",
				time.Millisecond * 500,
			},
			true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &task{
				f:        tt.fields.f,
				name:     tt.fields.name,
				duration: tt.fields.duration,
			}
			ctx, _ := context.WithTimeout(context.TODO(), time.Second*2)
			go func() {
				if err := t.Do(ctx); (err != nil) != tt.wantErr {
					t1.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			<-ctx.Done()

			if !(counter.Get() >= 3) || !(counter.Get() <= 4) {
				t1.Errorf("A different result was expected")
			}
		})
	}
}
