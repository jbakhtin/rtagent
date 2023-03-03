package memstorage

import (
	"context"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
)

type Snapshot struct {
	ToFile   *toFile
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
		ToFile:   newWriter,
	}, nil
}

func (s *Snapshot) Import(ctx context.Context) (map[string]models.Metric, bool) {
	list, err := s.FromFile.ReadList()
	if err != nil {
		return nil, false
	}

	return list, true
}

func (s *Snapshot) Exporting(ctx context.Context, cfg config.Config, metrics *map[string]models.Metric) {
	ticker := time.NewTicker(cfg.StoreInterval)

	// TODO: добавить условие для параметра
	for {
		select {
		case <-ctx.Done():
			err := s.ToFile.WriteList(metrics)
			if err != nil {
				return
			}
			_ = s.ToFile.Close()
			_ = s.FromFile.Close()
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
