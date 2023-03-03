package memstorage

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
)

type Reader interface {
	Read() (*models.Metric, error)
	Close() error
}

type fromFile struct {
	file   *os.File
	reader *bufio.Reader
}

func NewReader(cfg config.Config) (*fromFile, error) {
	var file *os.File
	var err error

	// открываем файл для чтения
	file, err = os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
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

	return &fromFile{
		file:   file,
		reader: bufio.NewReader(file),
	}, nil
}

func (ff *fromFile) ReadList() (map[string]models.Metric, error) {
	data, err := ff.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]models.Metric, 20)
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return metrics, err
	}

	return metrics, nil
}

func (ff *fromFile) Close() error {
	return ff.file.Close()
}
