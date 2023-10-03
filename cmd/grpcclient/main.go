package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "github.com/jbakhtin/rtagent/proto/metric"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := pb.NewMetricsClient(conn)

	// функция, в которой будем отправлять сообщения

	metricRequest := pb.UpdateMetricRequest{
		Metric: &pb.Metric{
			Type: pb.Metric_gauge,
			Value: 12,
			Key: "test",
		},
	}

	c.UpdateMetric(context.TODO(), &metricRequest)
}
