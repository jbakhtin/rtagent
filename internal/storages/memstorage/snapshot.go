package memstorage

import (
	"context"
	"fmt"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"time"
)

type Snapshot struct {
	ToFile *toFile
	FromFile *fromFile

	metrics *map[string]models.Metricer
}

func NewSnapshot(ctx context.Context, cfg config.Config) (*Snapshot, error) {
	newReader, err := NewReader(cfg)
	if err != nil {
		return nil, err
	}

	newWriter, err := NewWriter(cfg)
	if err != nil {
		return nil, err
	}

	return &Snapshot{
		FromFile: newReader,
		ToFile: newWriter,
	}, nil
}

func (s *Snapshot) Import(ctx context.Context) (map[string]models.Metric, error) {
	var metrics map[string]models.Metric
	metrics, err := s.FromFile.ReadList()
	if err != nil {
		fmt.Println(err)
	}

	return metrics, nil
}

func (s *Snapshot) Exporting(ctx context.Context, cfg config.Config, metrics *map[string]models.Metric) {
	ticker := time.NewTicker(cfg.StoreInterval)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("test")
			err := s.ToFile.WriteList(metrics)
			if err != nil {
				return
			}
			return
		case <-ticker.C:
			err := s.ToFile.WriteList(metrics)
			if err != nil {
				return
			}
		}
	}
}


func (s *Snapshot) Export(ctx context.Context, metrics map[string]models.Metric) error {
	err := s.ToFile.WriteList(&metrics)
	if err != nil {
		return err
	}

	return nil
}