package memstorage

import (
	"bufio"
	"encoding/json"
	"github.com/jbakhtin/rtagent/internal/config"
	"github.com/jbakhtin/rtagent/internal/models"
	"log"
	"os"
)

type Writer interface {
	Write(event *models.Metric) // для записи события
	Close() error            // для закрытия ресурса (файла)
}

type toFile struct {
	file *os.File // файл для записи
	writer *bufio.Writer
}

func NewWriter(cfg config.Config) (*toFile, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE, 0777)
	if os.IsNotExist(err) {
		err = os.Mkdir("tmp", 0777)
		if err != nil {
			log.Fatal(err)
		}

		file, err = os.OpenFile(cfg.StoreFile, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &toFile{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil
}

func (tf *toFile) WriteList(event *map[string]models.Metric) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := tf.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	err = tf.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = tf.file.Seek(0, 0)
	if err != nil {
		return err
	}
	if err = tf.writer.WriteByte('\n'); err != nil {
		return err
	}
	// записываем буфер в файл
	return tf.writer.Flush()
}

func (tf *toFile) Close() error {
	// закрываем файл
	return tf.file.Close()
	return nil
}