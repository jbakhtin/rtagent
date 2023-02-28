package memstorage

import (
	"bufio"
	"encoding/json"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"os"
)

type Reader interface {
	Read() (*models.Metric, error) // для чтения события
	Close() error               // для закрытия ресурса (файла)
}

type fromFile struct {
	file *os.File // файл для записи
	reader *bufio.Reader
}

func NewReader(cfg config.Config) (*fromFile, error) {

	// открываем файл для чтения
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &fromFile{
		file: file,
		reader: bufio.NewReader(file),
	}, nil
}

func (ff *fromFile) ReadList() (map[string]models.Metric, error) {
	// читаем данные до символа переноса строки
	data, err := ff.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	event := make(map[string]models.Metric, 10)
	// преобразуем данные из JSON-представления в структуру
	err = json.Unmarshal(data, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}

func (ff *fromFile) Close() error {
	// закрываем файл
	//return ff.file.Close()
	return nil
}