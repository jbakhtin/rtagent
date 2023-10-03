package grpcserver

import (
	"context"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/internal/types"
	pb "github.com/jbakhtin/rtagent/proto/metric"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	grpc.Server
	pb.UnimplementedMetricsServer

	Repository storage.MetricRepository
}

func New(cfg config.Config, repository storage.MetricRepository) (*Server, error) { //ToDo: need remove confog
	s := &Server{
		Repository: repository,
		Server:     *grpc.NewServer(),
	}

	pb.RegisterMetricsServer(s, s)
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
	case pb.Metric_counter:
		metric, err = models.NewCounter(types.CounterType, request.Metric.Key, strconv.Itoa(int(request.Metric.Delta)))
	case pb.Metric_gauge:
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
