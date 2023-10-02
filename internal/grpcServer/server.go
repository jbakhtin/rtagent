package grpcServer

import (
	"context"
	"fmt"
	pb "github.com/jbakhtin/rtagent/proto/metric"
	"sync"
)

// Server поддерживает все необходимые методы сервера.
type Server struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedMetricsServer

	// используем sync.Map для хранения пользователей
	metrics sync.Map
}

func (s *Server) UpdateMetric(ctx context.Context, metric *pb.UpdateMetricRequest) (*pb.UpdateMetricResponse, error) {
	var response pb.UpdateMetricResponse

	fmt.Println("метрика добавлена")


	return &response, nil
}