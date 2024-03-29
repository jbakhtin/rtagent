package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/go-faster/errors"
	"io"
	"log"
	"os"
	"time"

	"github.com/jbakhtin/rtagent/internal/models"
	handlerModels "github.com/jbakhtin/rtagent/internal/server/models"
	"github.com/jbakhtin/rtagent/internal/types"

	"github.com/jbakhtin/rtagent/internal/storage/memstorage"

	"github.com/jbakhtin/rtagent/internal/config"
)

// FileStorage является оберткой над MemStorage и вынесен в отдельный пакет, как полноценное хранилище
type FileStorage struct {
	memstorage.MemStorage
}

func (fs *FileStorage) Start(ctx context.Context, cfg config.Config) error {
	ticker := time.NewTicker(cfg.StoreInterval)
	defer ticker.Stop()

	if cfg.Restore {
		err := fs.Restore(ctx, cfg)
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := fs.Backup(ctx, cfg)
				if err != nil {
					log.Println(err) //ToDo need refactoring, need add channels
				}
				return

			case <-ticker.C:
				err := fs.Backup(ctx, cfg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	return nil
}

func (fs *FileStorage) Backup(ctx context.Context, cfg config.Config) error {
	file, err := fs.openFile(cfg, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {

		return errors.Wrap(err, "open file")
	}
	defer file.Close()

	metrics, err := fs.GetAll()
	if err != nil {
		return errors.Wrap(err, "mem storage")
	}

	var JSONMetrics []handlerModels.Metrics
	var JSONMetric handlerModels.Metrics

	for _, v := range metrics {
		JSONMetric, err = v.ToJSON([]byte(cfg.KeyApp))
		if err != nil {
			return errors.Wrap(err, "metric to json")
		}
		JSONMetrics = append(JSONMetrics, JSONMetric)
	}

	data, err := json.Marshal(JSONMetrics)
	if err != nil {
		return errors.Wrap(err, "marshal json")
	}

	writer := bufio.NewWriter(file)
	if _, err = writer.Write(data); err != nil {
		return errors.Wrap(err, "write to buffer")
	}

	err = file.Truncate(int64(0))
	if err != nil {
		return errors.Wrap(err, "truncate file")
	}

	if _, err = file.Seek(0, 0); err != nil {
		return errors.Wrap(err, "seek file")
	}

	if err = writer.WriteByte('\n'); err != nil {
		return errors.Wrap(err, "write bytes to file with \\\n")
	}

	return writer.Flush()
}

func (fs *FileStorage) Restore(ctx context.Context, cfg config.Config) error {
	file, err := fs.openFile(cfg, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return closeErr
		}

		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			return errors.New("test")
		}

		return nil
	}

	JSONMetrics := make([]handlerModels.Metrics, 0)
	err = json.Unmarshal(data, &JSONMetrics)
	if err != nil {
		return err
	}

	for _, JSONMetric := range JSONMetrics {
		switch JSONMetric.MType {
		case types.GaugeType:
			fs.Items[JSONMetric.MKey] = models.Gauge{
				Description: models.Description{
					MKey:  JSONMetric.MKey,
					MType: JSONMetric.MType,
				},
				MValue: *JSONMetric.Value,
			}
		case types.CounterType:
			fs.Items[JSONMetric.MKey] = models.Counter{
				Description: models.Description{
					MKey:  JSONMetric.MKey,
					MType: JSONMetric.MType,
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
		//fs.Logger.Info(err.Error())

		//fs.Logger.Info("try to make dir 'tmp'")
		err = os.Mkdir("./tmp", perm)
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
