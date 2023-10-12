package grpcserver

import (
	"context"
	"fmt"
	"github.com/go-faster/errors"
	pb "github.com/jbakhtin/rtagent/gen/go/metric/v1"
	"github.com/jbakhtin/rtagent/internal/config"
	trustedsubnets "github.com/jbakhtin/rtagent/internal/grpcserver/interceptors"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/storage"
	"github.com/jbakhtin/rtagent/internal/types"
	"github.com/jbakhtin/rtagent/pkg/hasher"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	grpc.Server
	pb.UnimplementedMetricsServiceServer

	Repository storage.MetricRepository
	cfg        config.Config
}

func New(cfg config.Config, repository storage.MetricRepository) (*Server, error) {
	trustedSubnets := trustedsubnets.TrustedSubnet{
		Subnets: cfg.TrustedSubnet,
	}

	s := &Server{
		Repository: repository,
		Server:     *grpc.NewServer(grpc.UnaryInterceptor(trustedSubnets.TrustedSubnetsInterceptor)),
		cfg:        cfg,
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
	var hash string

	switch request.Metric.Type {
	case pb.Type_TYPE_COUNTER:
		metric, err = models.NewCounter(types.CounterType, request.Metric.Key, strconv.Itoa(int(request.Metric.Delta)))
		if s.cfg.GetKeyApp() != "" {
			hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%v", request.Metric.Key, request.Metric.Type, request.Metric.Delta), []byte(s.cfg.GetKeyApp()))
			if err != nil {
				return nil, errors.Wrap(err, "calc hash")
			}

			if hash != request.Metric.Hash {
				return nil, errors.Wrap(err, "hash not equal")
			}
		}
	case pb.Type_TYPE_GAUGE:
		metric, err = models.NewGauge(types.GaugeType, request.Metric.Key, strconv.FormatFloat(float64(request.Metric.Value), 'E', -1, 32))
		if s.cfg.GetKeyApp() != "" {
			hash, err = hasher.CalcHash(fmt.Sprintf("%s:%s:%v", request.Metric.Key, request.Metric.Type, request.Metric.Value), []byte(s.cfg.GetKeyApp()))
			if err != nil {
				return nil, errors.Wrap(err, "calc hash")
			}

			if hash != request.Metric.Hash {
				return nil, errors.Wrap(err, "hash not equal")
			}
		}
	default:
		response.Error = "metric typ not valid"
		return &response, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = s.Repository.Set(metric) // ToDo: need pass the context into
	if err != nil {
		return nil, err
	}

	return &response, nil
}
