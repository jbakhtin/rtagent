package postgres

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	_ "github.com/jackc/pgx/v5"
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

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DBStorage struct {
	DatabaseDSN string
	Driver      string
	Logger      *zap.Logger
}

func New(driver string, cfg config.Config) (DBStorage, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return DBStorage{}, err
	}

	db, err := sql.Open(driver, cfg.DatabaseDSN) // TODO: можно ли открыть соединение один раз при старте приложения и закрыть при остановке?
	if err != nil {
		return DBStorage{}, err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return DBStorage{}, err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return DBStorage{}, err
	}

	return DBStorage{
		DatabaseDSN: cfg.DatabaseDSN,
		Logger:      logger,
		Driver:      driver,
	}, nil
}

func (dbs *DBStorage) Set(metric models.Metricer) (models.Metricer, error) {
	db, err := sql.Open(dbs.Driver, dbs.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var newMetric models.Metricer
	switch m := metric.(type) {
	case models.Gauge:
		var newM models.Gauge
		err = db.QueryRow(insertGauge, m.MKey, m.MType, m.MValue).Scan(&newM.MKey, &newM.MType, &newM.MValue)
		if err != nil {
			return nil, err
		}
		newMetric = newM
	case models.Counter:
		var newM models.Counter
		err = db.QueryRow(insertCounter, m.MKey, m.MType, m.MValue).Scan(&newM.MKey, &newM.MType, &newM.MValue)
		if err != nil {
			return nil, err
		}
		newMetric = newM
	}

	return newMetric, nil
}

func (dbs *DBStorage) Get(key string) (models.Metricer, error) {
	db, err := sql.Open(dbs.Driver, dbs.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var id string
	var mType string
	var delta *types.Counter
	var value *types.Gauge
	err = db.QueryRow(getByID, key).Scan(&id, &mType, &delta, &value)
	if err != nil {
		dbs.Logger.Info(err.Error())
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

func (dbs *DBStorage) GetAll() (map[string]models.Metricer, error) {
	db, err := sql.Open(dbs.Driver, dbs.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(getAll) // TODO: need limit
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
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

func (dbs *DBStorage) SetBatch(ctx context.Context, metrics []models.Metricer) ([]models.Metricer, error){
	db, err := sql.Open("pgx", dbs.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	stmtGauge, err := tx.PrepareContext(ctx, insertGauge)
	if err != nil {
		return nil, err
	}
	defer stmtGauge.Close()

	stmtCounter, err := tx.PrepareContext(ctx, insertCounter)
	if err != nil {
		return nil, err
	}
	defer stmtCounter.Close()

	for _, v := range metrics {
		switch metric := v.(type) {
		case models.Gauge:
			if _, err = stmtGauge.ExecContext(ctx, metric.MKey, metric.MType, metric.MValue); err != nil {
				return nil, err
			}
		case models.Counter:
			if _, err = stmtCounter.ExecContext(ctx, metric.MKey, metric.MType, metric.MValue); err != nil {
				return nil, err
			}
		}
	}
	tx.Commit()

	return metrics, nil
}
