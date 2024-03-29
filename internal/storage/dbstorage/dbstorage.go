package dbstorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5"
	"github.com/jbakhtin/rtagent/internal/models"
	"github.com/jbakhtin/rtagent/internal/types"

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
var EmbedMigrations embed.FS

type DBStorage struct {
	*sql.DB
}

func (dbs *DBStorage) Set(metric models.Metricer) (models.Metricer, error) {
	switch m := metric.(type) {
	case models.Gauge:
		var metricSaved models.Gauge
		err := dbs.QueryRow(insertGauge, m.MKey, m.MType, m.MValue). // ToDo: Need to use QueryContext
										Scan(&metricSaved.MKey, &metricSaved.MType, &metricSaved.MValue)
		if err != nil {
			return nil, err
		}
		return metricSaved, nil
	case models.Counter:
		var metricSaved models.Counter
		err := dbs.QueryRow(insertCounter, m.MKey, m.MType, m.MValue).
			Scan(&metricSaved.MKey, &metricSaved.MType, &metricSaved.MValue)
		if err != nil {
			return nil, err
		}
		return metricSaved, nil
	}

	return nil, errors.New("type not recognized")
}

func (dbs *DBStorage) Get(key string) (models.Metricer, error) {
	var id string
	var mType string
	var delta *types.Counter
	var value *types.Gauge
	err := dbs.QueryRow(getByID, key).Scan(&id, &mType, &delta, &value)
	if err != nil {
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
	rows, err := dbs.Query(getAll) // TODO: need limit
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

func (dbs *DBStorage) SetBatch(metrics []models.Metricer) ([]models.Metricer, error) {
	tx, err := dbs.Begin()
	if err != nil {
		return nil, err
	}

	defer func(tx *sql.Tx) {
		if tempErr := tx.Rollback(); tempErr != nil {
			err = tempErr
		}
	}(tx)

	ctx := context.TODO()

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
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func (dbs *DBStorage) TestPing() error {
	return dbs.Ping()
}
