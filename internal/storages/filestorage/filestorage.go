package filestorage

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/jbakhtin/rtagent/internal/storages/memstorage"
	"io"
	"log"
	"os"
	"time"

	"github.com/jbakhtin/rtagent/internal/config"
)

type FileStorage struct {
	memstorage.MemStorage
}

func New(ctx context.Context, cfg config.Config) (FileStorage, error) { // TODO: определить интерфейс
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

	go func () {
		for {
			select {
			case <-ctx.Done():
				err := fs.Write(ctx, cfg)
				if err != nil {
					//return err
				}
				//return nil
			case <-ticker.C:
				err := fs.Write(ctx, cfg)
				if err != nil {
					//return err
				}
			}
		}
	}()

	return nil
}

func (fs *FileStorage) Write(ctx context.Context, cfg config.Config) error {
	file, err := os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if os.IsNotExist(err) {
		// TODO: файл создается только если присутсвует указанная директория, можно ли как то по дургому?
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
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
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	if os.IsNotExist(err) {
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err) // TODO: выкинуть ошибку через канал
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

	err = json.Unmarshal(data, &fs.Items)
	if err != nil {
		return err
	}

	return nil
}




