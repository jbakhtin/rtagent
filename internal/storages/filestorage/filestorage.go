package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
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

func New(ctx context.Context, cfg config.Config) (FileStorage, error) {
	memStorage, err := memstorage.NewMemStorage(ctx, cfg)
	if err != nil {
		return FileStorage{}, err
	}

	return FileStorage{
		MemStorage: memStorage,
	}, nil
}

func (fs *FileStorage) Start(ctx context.Context, cfg config.Config) error {
	ticker := time.NewTicker(cfg.StoreInterval)

	err := fs.Read(ctx, cfg)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := fs.Write(ctx, cfg)
				if err != nil {
					fs.Logger.Error(err.Error())
				} else {
					fs.Logger.Info("the data is saved to the disk")
				}

			case <-ticker.C:
				err := fs.Write(ctx, cfg)
				if err != nil {
					fs.Logger.Error(err.Error())
				}
			}
		}
	}()

	return nil
}

func (fs *FileStorage) Write(ctx context.Context, cfg config.Config) error {
	file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if os.IsNotExist(err) {
		fs.Logger.Info(err.Error())

		fs.Logger.Info("try to make dir 'tmp'")
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			return err
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(fs.Items)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		return err
	}

	// Удаляем содержимое файла перед перезаписью
	// TODO: Не удалось решить данную проблему дргим способом, например ...|os.O_TRUNC
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

func (fs *FileStorage) Read(ctx context.Context, cfg config.Config) error {
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0644)
	if os.IsNotExist(err) {
		fs.Logger.Info(err.Error())

		fs.Logger.Info("try to make dir 'tmp'")
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			return err
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
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

	metrics, err := fs.GetAll()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return err
	}

	return nil
}
