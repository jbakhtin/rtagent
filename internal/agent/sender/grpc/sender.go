package grpc

import (
	"context"
	pb "github.com/jbakhtin/rtagent/gen/go/metric/v1"
	"github.com/jbakhtin/rtagent/internal/agent/sender/grpc/interceptors"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strconv"
)

type Configer interface {
	GetServerAddress() string
	GetKeyApp() string
	GetCryptoKey() string
	GetTrustedSubnet() string
	GetGRPCServerAddress() string
}

type ReportFunction func() string

type grpcSender struct {
	cfg  Configer
	conn *grpc.ClientConn
}

func New(cfg Configer) (*grpcSender, error) {
	xRealIP := xrealip.XRealIP{
		IPs: cfg.GetTrustedSubnet(),
	}

	conn, err := grpc.Dial(
		cfg.GetGRPCServerAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			xRealIP.SetXRealIPInterceptor,
		),
	)
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

	var pbMetric *pb.Metric

	switch v := value.(type) {
	case types.Counter:
		metricModel, err := models.NewCounter("counter", key, strconv.FormatInt(int64(v), 10))
		if err != nil {
			return err
		}
		pbMetric, err = metricModel.ToGRPC(r.cfg.GetKeyApp())
		if err != nil {
			return err
		}
	case types.Gauge:
		metricModel, err := models.NewGauge("gauge", key, strconv.FormatFloat(float64(v), 'E', -1, 64))
		if err != nil {
			return err
		}
		pbMetric, err = metricModel.ToGRPC(r.cfg.GetKeyApp())
		if err != nil {
			return err
		}
	}

	metricRequest := pb.UpdateMetricRequest{
		Metric: pbMetric,
	}

	//ToDo: need log response error
	_, err := c.UpdateMetric(context.TODO(), &metricRequest)
	if err != nil {
		return err
	}

	return nil
}
