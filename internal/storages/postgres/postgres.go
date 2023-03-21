package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	insertGauge = `
		INSERT INTO metrics (id, type, value)
		VALUES($1, $2, $3)
		ON CONFLICT (id, type) 
		DO UPDATE SET value = $3
		RETURNING id, type, value;
	`
	insertCounter = `
		INSERT INTO metrics (id, type, delta)
		VALUES($1, $2, $3)
		ON CONFLICT (id, type) 
		DO UPDATE SET delta = (metrics.delta + ($3))
		RETURNING id, type, delta;
	`

	getByID = `
		SELECT id, type, delta, value FROM metrics
		WHERE metrics.id = $1
	`

	getAll = `
		SELECT id, type, delta, value FROM metrics
	`
)

//go:embed migrations/20230319143358_create_metrics_table.sql
var embedMigrations embed.FS

type MemStorage struct {
	DatabaseDSN string
	Logger      *zap.Logger
}

func New(cfg config.Config) (MemStorage, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return MemStorage{}, err
	}

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return MemStorage{}, err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return MemStorage{}, err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return MemStorage{}, err
	}

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

	var newMetric models.Metricer
	switch m := metric.(type) {
	case models.Gauge:
		var newM models.Gauge
		err = conn.QueryRow(context.Background(), insertGauge, m.MKey, m.MType, m.MValue).Scan(&newM.MKey, &newM.MType, &newM.MValue)
		if err != nil {
			return nil, err
		}
		newMetric = newM
	case models.Counter:
		var newM models.Counter
		err = conn.QueryRow(context.Background(), insertCounter, m.MKey, m.MType, m.MValue).Scan(&newM.MKey, &newM.MType, &newM.MValue)
		if err != nil {
			return nil, err
		}
		newMetric = newM
	}

	return newMetric, nil
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
	err = conn.QueryRow(context.Background(), getByID, key).Scan(&id, &mType, &delta, &value)
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
	if err != nil {
		return nil, err
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
			return nil, err
		}

		switch mType {
		case types.GaugeType:
			metric, err = models.NewGauge(mType, id, fmt.Sprintf("%v", *value))
		case types.CounterType:
			metric, err = models.NewCounter(mType, id, fmt.Sprintf("%v", *delta))
		}
		if err != nil {
			return nil, err
		}

		metrics[metric.Key()] = metric
	}

	return metrics, nil
}

func (pg *MemStorage) SetBatch(ctx context.Context, metrics []models.Metricer) ([]models.Metricer, error){

	db, err := sql.Open("pgx", pg.DatabaseDSN)

	// шаг 1 — объявляем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	// шаг 1.1 — если возникает ошибка, откатываем изменения
	defer tx.Rollback()

	// шаг 2 — готовим инструкцию
	stmtGauge, err := tx.PrepareContext(ctx, insertGauge)
	stmtCounter, err := tx.PrepareContext(ctx, insertCounter)
	if err != nil {
		return nil, err
	}
	// шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
	defer stmtGauge.Close()
	defer stmtCounter.Close()

	for _, v := range metrics {
		// шаг 3 — указываем, что каждое видео будет добавлено в транзакцию
		switch metric := v.(type) {
		case models.Gauge:
			if _, err = stmtGauge.ExecContext(ctx, metric.MKey, metric.MType, metric.MValue); err != nil {
				return nil, err
			}
		case models.Counter:
			fmt.Println(metric)
			if _, err = stmtCounter.ExecContext(ctx, metric.MKey, metric.MType, metric.MValue); err != nil {
				return nil, err
			}
		}
	}
	// шаг 4 — сохраняем изменения
	tx.Commit()

	return metrics, nil
}
