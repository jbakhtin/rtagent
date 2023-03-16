package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/jbakhtin/rtagent/internal/models"
	models2 "github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"
	"io"
	"os"
	"time"

	"github.com/jbakhtin/rtagent/internal/storages/memstorage"

	"github.com/jbakhtin/rtagent/internal/config"
)

//FileStorage является оберткой над MemStorage и вынесен в отдельный пакет, как полноценное хранилище
type FileStorage struct {
	memstorage.MemStorage
}

func New(cfg config.Config) (FileStorage, error) {
	memStorage, err := memstorage.NewMemStorage(cfg)
	if err != nil {
		return FileStorage{}, err
	}

	return FileStorage{
		MemStorage: memStorage,
	}, nil
}

func (fs *FileStorage) Start(ctx context.Context, cfg config.Config) error {
	ticker := time.NewTicker(cfg.StoreInterval)

	err := fs.Restore(ctx, cfg)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := fs.Backup(ctx, cfg)
				if err != nil {
					fs.Logger.Error(err.Error())
				} else {
					fs.Logger.Info("the data is saved to the disk")
				}
				fs.Logger.Info("stop file storage loop")
				return

			case <-ticker.C:
				err := fs.Backup(ctx, cfg)
				if err != nil {
					fs.Logger.Error(err.Error())
				}
			}
		}
	}()

	return nil
}

func (fs *FileStorage) Backup(ctx context.Context, cfg config.Config) error {
	file, err := fs.openFile(cfg, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	metrics, err := fs.GetAll()
	if err != nil {
		return err
	}

	var JSONMetrics []models2.Metrics

	for _, v := range metrics {
		JSONMetrics = append(JSONMetrics, v.ToJSON())
	}

	data, err := json.Marshal(JSONMetrics)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}

	err = file.Truncate(0)
	if err != nil {
		return err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return err
	}

	if err = writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()
}

func (fs *FileStorage) Restore(ctx context.Context, cfg config.Config) error {
	file, err := fs.openFile(cfg, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			return err
		}

		return nil
	}

	JSONMetrics := make([]models2.Metrics, 0)
	err = json.Unmarshal(data, &JSONMetrics)
	if err != nil {
		return err
	}

	for _, JSONMetric := range JSONMetrics {
		switch JSONMetric.MType {
		case types.GaugeType:
			fs.Items[JSONMetric.MKey] = models.Gauge{
				Description: models.Description{
					MKey:  JSONMetric.MType,
					MType: JSONMetric.MKey,
				},
				MValue: *JSONMetric.Value,
			}
		case types.CounterType:
			fs.Items[JSONMetric.MKey] = models.Counter{
				Description: models.Description{
					MKey:  JSONMetric.MType,
					MType: JSONMetric.MKey,
				},
				MValue: *JSONMetric.Delta,
			}
		}
	}

	return nil
}

func (fs *FileStorage) openFile(cfg config.Config, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(cfg.StoreFile, flag, perm)
	if os.IsNotExist(err) {
		fs.Logger.Info(err.Error())

		fs.Logger.Info("try to make dir 'tmp'")
		err = os.Mkdir("tmp", perm)
		if err != nil {
			return nil, err
		}

		file, err = os.OpenFile(cfg.StoreFile, flag, perm)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}
