package grpcServer

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/internal/types"
	pb "github.com/jbakhtin/rtagent/proto/metric"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

// Server поддерживает все необходимые методы сервера.
type Server struct {
	grpc.Server
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	Repository storage.MetricRepository
	config     config.Config
}

func New(cfg config.Config, repository storage.MetricRepository) *Server {

	s := &Server{
		Repository: repository,
		Server: *grpc.NewServer(),
	}

	pb.RegisterMetricsServer(s, s)

	return s
}

func (s *Server) Run(ctx context.Context) (err error) {
	// определяем порт для сервера
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		return
	}

	fmt.Println("Се" +
		"" +
		"рвер gRPC начал работу")
	// получаем запрос gRPC
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

	fmt.Println("test")

	switch request.Metric.Type {
	case pb.Metric_counter:
		metric, err = models.NewCounter(types.CounterType, request.Metric.Key, strconv.Itoa(int(request.Metric.Delta)))
	case pb.Metric_gauge:
		metric, err = models.NewGauge(types.GaugeType, request.Metric.Key, strconv.FormatFloat(float64(request.Metric.Value), 'E', -1, 32))
	default:
		response.Error = "metric typ not valid"
		return &response, nil
	}

	_, err = s.Repository.Set(metric)
	if err != nil {
		return nil, err
	}

	return &response, nil
}