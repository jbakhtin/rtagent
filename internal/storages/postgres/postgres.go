package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"go.uber.org/zap"
)

const (
	insert = `
		INSERT INTO metrics (id, type, delta, value)
		VALUES($1, $2, $3, $4)
		ON CONFLICT (id, type) 
		DO UPDATE SET delta = $3, value = $4
		WHERE metrics.id = $1
		RETURNING id;
	`

	getById = `
		SELECT id, type, delta, value FROM metrics
		WHERE metrics.id = $1
	`

	getAll = `
		SELECT id, type, delta, value FROM metrics
	`
)

type MemStorage struct {
	DatabaseDSN string
	Logger      *zap.Logger
}

func New(cfg config.Config) (MemStorage, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return MemStorage{}, err
	}

	logger.Info(cfg.DatabaseDSN)
	logger.Info("test")

	return MemStorage{
		DatabaseDSN: cfg.DatabaseDSN,
		Logger:      logger,
	}, nil
}

func (ms *MemStorage) Set(metric models.Metricer) (models.Metricer, error) {
	conn, err := pgx.Connect(context.Background(), ms.DatabaseDSN)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}
	defer conn.Close(context.Background())

	var id string
	JSONMetric, err := metric.ToJSON(nil)
	err = conn.QueryRow(context.Background(), insert, JSONMetric.MKey, JSONMetric.MType, JSONMetric.Delta, JSONMetric.Value).Scan(&id)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}

	return metric, nil
}

func (ms *MemStorage) Get(key string) (models.Metricer, error) {
	conn, err := pgx.Connect(context.Background(), ms.DatabaseDSN)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}
	defer conn.Close(context.Background())

	var id string
	var mType string
	var delta *types.Counter
	var value *types.Gauge
	err = conn.QueryRow(context.Background(), getById, key).Scan(&id, &mType, &delta, &value)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}

	var metric models.Metricer
	switch mType {
	case types.GaugeType:
		metric, err = models.NewGauge(mType, id, fmt.Sprintf("%v", *value))
	case types.CounterType:
		metric, err = models.NewCounter(mType, id, fmt.Sprintf("%v", *delta))
	}

	return metric, nil
}

func (ms *MemStorage) GetAll() (map[string]models.Metricer, error) {
	conn, err := pgx.Connect(context.Background(), ms.DatabaseDSN)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), getAll)
	if err != nil {
		ms.Logger.Info(err.Error())
		return nil, err
	}

	var id string
	var mType string
	var delta *types.Counter
	var value *types.Gauge
	var metric models.Metricer
	metrics := make(map[string]models.Metricer)

	for rows.Next() {
		err = rows.Scan(&id, &mType, &delta, &value)
		if err != nil {
			ms.Logger.Info(err.Error())
			return nil, err
		}

		switch mType {
		case types.GaugeType:
			metric, err = models.NewGauge(mType, id, fmt.Sprintf("%v", *value))
		case types.CounterType:
			metric, err = models.NewCounter(mType, id, fmt.Sprintf("%v", *delta))
		}

		metrics[metric.Key()] = metric
	}

	return metrics, nil
}
