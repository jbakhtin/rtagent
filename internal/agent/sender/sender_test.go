package sender

import (
	"github.com/jbakhtin/rtagent/internal/agent/sender/http"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/types"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	cfg, err := config.NewConfigBuilder().WithAllFromEnv().Build()
	if err != nil {
		t.Errorf("Config builder: %v", err)
	}

	type args struct {
		cfg http.Configer
	}
	tests := []struct {
		name string
		args args
		want *http.HttpSender
	}{
		{
			"New sender",
			args{
				cfg,
			},
			&http.HttpSender{
				cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := http.New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sender_Send(t *testing.T) {
	cfg, err := config.NewConfigBuilder().WithAllFromFlagsA().WithAllFromEnv().Build()
	if err != nil {
		t.Errorf("Config builder: %v", err)
	}

	type fields struct {
		cfg http.Configer
	}
	type args struct {
		key   string
		value types.Metricer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Send to not started server",
			fields{
				cfg,
			},
			args{
				"TestMetric",
				types.Counter(10),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.HttpSender{
				Cfg: tt.fields.cfg,
			}
			if err := r.Send(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
