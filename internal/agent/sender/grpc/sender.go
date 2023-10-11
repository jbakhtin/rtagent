package grpc

import (
	"context"
	pb "github.com/jbakhtin/rtagent/gen/go/metric/v1"
	"github.com/jbakhtin/rtagent/internal/types"
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
	cfg  Configer
	conn *grpc.ClientConn
}

func New(cfg Configer) (*grpcSender, error) {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &grpcSender{
		cfg:  cfg,
		conn: conn,
	}, nil
}

func (r *grpcSender) Send(key string, value types.Metricer) error {
	c := pb.NewMetricsServiceClient(r.conn)

	metric := &pb.Metric{
		Key: key,
	}
	switch v := value.(type) {
	case types.Counter:
		metric.Delta = uint64(v)
		metric.Type = pb.Type_TYPE_COUNTER
	case types.Gauge:
		metric.Value = float32(v)
		metric.Type = pb.Type_TYPE_GAUGE
	}
	metric.Hash = "test" //ToDo: need implement hash calc

	metricRequest := pb.UpdateMetricRequest{
		Metric: metric,
	}

	//ToDo: need log response error
	_, err := c.UpdateMetric(context.TODO(), &metricRequest)
	if err != nil {
		return err
	}

	return nil
}
