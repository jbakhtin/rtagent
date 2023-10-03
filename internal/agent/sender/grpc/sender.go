package grpc

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/types"
	pb "github.com/jbakhtin/rtagent/proto/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Configer interface {
	GetServerAddress() string
	GetKeyApp() string
	GetCryptoKey() string
	GetTrustedSubnet() string
}

type ReportFunction func() string

type grpcSender struct {
	cfg Configer
	conn *grpc.ClientConn
}

func New(cfg Configer) *grpcSender {
	conn, _ := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))

	return &grpcSender{
		cfg: cfg,
		conn: conn,
	}
}

func (r *grpcSender) Send(key string, value types.Metricer) error {
	c := pb.NewMetricsClient(r.conn)

	// функция, в которой будем отправлять сообщения

	metric := &pb.Metric{
		Key: key,
	}
	switch v := value.(type) {
	case types.Counter:
		metric.Delta = uint64(v)
		metric.Type = pb.Metric_counter
	case types.Gauge:
		metric.Value = float32(v)
		metric.Type = pb.Metric_gauge
	}
	metric.Hash = "test"

	metricRequest := pb.UpdateMetricRequest{
		Metric: metric,
	}

	c.UpdateMetric(context.TODO(), &metricRequest)

	return nil
}
