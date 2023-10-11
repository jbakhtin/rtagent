package grpcserver

import (
	"context"
	pb "github.com/jbakhtin/rtagent/gen/go/metric/v1"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/internal/types"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	grpc.Server
	pb.UnimplementedMetricsServiceServer

	Repository storage.MetricRepository
}

func New(cfg config.Config, repository storage.MetricRepository) (*Server, error) { //ToDo: need remove confog
	s := &Server{
		Repository: repository,
		Server:     *grpc.NewServer(),
	}

	pb.RegisterMetricsServiceServer(s, s)
	//ToDo: need implement cors

	return s, nil
}

func (s *Server) Run(ctx context.Context) (err error) {
	listen, err := net.Listen("tcp", ":3200") //ToDo: need move grpc server address to config
	if err != nil {
		return
	}

	go func() {
		if err = s.Serve(listen); err != nil {
			return
		}
	}()

	return
}

func (s *Server) UpdateMetric(ctx context.Context, request *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse
	var metric models.Metricer
	var err error

	switch request.Metric.Type {
	case pb.Metric_TYPE_COUNTER:
		metric, err = models.NewCounter(types.CounterType, request.Metric.Key, strconv.Itoa(int(request.Metric.Delta)))
	case pb.Metric_TYPE_GAUGE_UNSPECIFIED:
		metric, err = models.NewGauge(types.GaugeType, request.Metric.Key, strconv.FormatFloat(float64(request.Metric.Value), 'E', -1, 32))
	default:
		response.Error = "metric typ not valid"
		return &response, nil
	}
	if err != nil {
		return nil, err
	}

	//ToDo: need to check th hash

	_, err = s.Repository.Set(metric) // ToDo: need pass the context into
	if err != nil {
		return nil, err
	}

	return &response, nil
}
